package server

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/rilldata/rill/admin"
	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/server/auth"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/rilldata/rill/runtime/pkg/email"
	"github.com/rilldata/rill/runtime/pkg/observability"
	runtimeauth "github.com/rilldata/rill/runtime/server/auth"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const prodDeplTTL = 14 * 24 * time.Hour

// runtimeAccessTokenTTL is the validity duration of JWTs issued for runtime access when calling GetProject.
// This TTL is not used for tokens created for internal communication between the admin and runtime services.
const runtimeAccessTokenDefaultTTL = 30 * time.Minute

// runtimeAccessTokenEmbedTTL is the validation duration of JWTs issued for embedding.
// Since low-risk embed users might not implement refresh, it defaults to a high value of 24 hours.
// It can be overridden to a lower value when issued for high-risk embed users.
const runtimeAccessTokenEmbedTTL = 24 * time.Hour

func (s *Server) ListProjectsForOrganization(ctx context.Context, req *adminv1.ListProjectsForOrganizationRequest) (*adminv1.ListProjectsForOrganizationResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.OrganizationName),
	)

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.OrganizationName)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := unmarshalPageToken(req.PageToken)
	if err != nil {
		return nil, err
	}
	pageSize := validPageSize(req.PageSize)

	// If user has ManageProjects, return all projects
	claims := auth.GetClaims(ctx)
	var projs []*database.Project
	if claims.OrganizationPermissions(ctx, org.ID).ManageProjects {
		projs, err = s.admin.DB.FindProjectsForOrganization(ctx, org.ID, token.Val, pageSize)
	} else if claims.OwnerType() == auth.OwnerTypeUser {
		// Get projects the user is a (direct or group) member of (note: the user can be a member of a project in the org, without being a member of org - we call this an "outside member")
		// plus all public projects
		projs, err = s.admin.DB.FindProjectsForOrgAndUser(ctx, org.ID, claims.OwnerID(), token.Val, pageSize)
	} else {
		projs, err = s.admin.DB.FindPublicProjectsInOrganization(ctx, org.ID, token.Val, pageSize)
	}
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// If no projects are public, and user is not an outside member of any projects, the projsMap is empty.
	// If additionally, the user is not an org member, return permission denied (instead of an empty slice).
	if len(projs) == 0 && !claims.OrganizationPermissions(ctx, org.ID).ReadProjects {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read projects")
	}

	nextToken := ""
	if len(projs) >= pageSize {
		nextToken = marshalPageToken(projs[len(projs)-1].Name)
	}

	dtos := make([]*adminv1.Project, len(projs))
	for i, p := range projs {
		dtos[i] = s.projToDTO(p, org.Name)
	}

	return &adminv1.ListProjectsForOrganizationResponse{
		Projects:      dtos,
		NextPageToken: nextToken,
	}, nil
}

func (s *Server) GetProject(ctx context.Context, req *adminv1.GetProjectRequest) (*adminv1.GetProjectResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.OrganizationName),
		attribute.String("args.project", req.Name),
	)

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.OrganizationName)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	proj, err := s.admin.DB.FindProjectByName(ctx, req.OrganizationName, req.Name)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("project %q not found", req.Name))
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if proj.Public {
		permissions.ReadProject = true
		permissions.ReadProd = true
	}
	if claims.Superuser(ctx) {
		permissions.ReadProject = true
		permissions.ReadProd = true
		permissions.ReadProdStatus = true
		permissions.ReadDev = true
		permissions.ReadDevStatus = true
		permissions.ReadProjectMembers = true
	}

	if !permissions.ReadProject {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project")
	}

	if proj.ProdDeploymentID == nil || !permissions.ReadProd {
		return &adminv1.GetProjectResponse{
			Project:            s.projToDTO(proj, org.Name),
			ProjectPermissions: permissions,
		}, nil
	}

	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !permissions.ReadProdStatus {
		depl.StatusMessage = ""
	}

	var attr map[string]any
	if claims.OwnerType() == auth.OwnerTypeUser {
		attr, err = s.jwtAttributesForUser(ctx, claims.OwnerID(), proj.OrganizationID, permissions)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	ttlDuration := runtimeAccessTokenDefaultTTL
	if req.AccessTokenTtlSeconds != 0 {
		ttlDuration = time.Duration(req.AccessTokenTtlSeconds) * time.Second
	}

	jwt, err := s.issuer.NewToken(runtimeauth.TokenOptions{
		AudienceURL: depl.RuntimeAudience,
		Subject:     claims.OwnerID(),
		TTL:         ttlDuration,
		InstancePermissions: map[string][]runtimeauth.Permission{
			depl.RuntimeInstanceID: {
				// TODO: Remove ReadProfiling and ReadRepo (may require frontend changes)
				runtimeauth.ReadObjects,
				runtimeauth.ReadMetrics,
				runtimeauth.ReadProfiling,
				runtimeauth.ReadRepo,
				runtimeauth.ReadAPI,
			},
		},
		Attributes: attr,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not issue jwt: %s", err.Error())
	}

	s.admin.Used.Deployment(depl.ID)

	return &adminv1.GetProjectResponse{
		Project:            s.projToDTO(proj, org.Name),
		ProdDeployment:     deploymentToDTO(depl),
		Jwt:                jwt,
		ProjectPermissions: permissions,
	}, nil
}

func (s *Server) SearchProjectNames(ctx context.Context, req *adminv1.SearchProjectNamesRequest) (*adminv1.SearchProjectNamesResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.pattern", req.NamePattern),
		attribute.Int("args.annotations", len(req.Annotations)),
	)

	claims := auth.GetClaims(ctx)
	if !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "only superusers can search projects")
	}

	token, err := unmarshalPageToken(req.PageToken)
	if err != nil {
		return nil, err
	}
	pageSize := validPageSize(req.PageSize)

	var projectNames []string
	if req.Annotations != nil && len(req.Annotations) > 0 {
		// If an annotation is set to "*", we just check for key presence (instead of exact key-value match)
		var annotationKeys []string
		for k, v := range req.Annotations {
			if v == "*" {
				annotationKeys = append(annotationKeys, k)
				delete(req.Annotations, k)
			}
		}

		projectNames, err = s.admin.DB.FindProjectPathsByPatternAndAnnotations(ctx, req.NamePattern, token.Val, annotationKeys, req.Annotations, pageSize)
	} else {
		projectNames, err = s.admin.DB.FindProjectPathsByPattern(ctx, req.NamePattern, token.Val, pageSize)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	nextToken := ""
	if len(projectNames) >= pageSize {
		nextToken = marshalPageToken(projectNames[len(projectNames)-1])
	}

	return &adminv1.SearchProjectNamesResponse{
		Names:         projectNames,
		NextPageToken: nextToken,
	}, nil
}

func (s *Server) CreateProject(ctx context.Context, req *adminv1.CreateProjectRequest) (*adminv1.CreateProjectResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.OrganizationName),
		attribute.String("args.project", req.Name),
		attribute.String("args.description", req.Description),
		attribute.Bool("args.public", req.Public),
		attribute.String("args.provisioner", req.Provisioner),
		attribute.String("args.prod_version", req.ProdVersion),
		attribute.String("args.prod_olap_driver", req.ProdOlapDriver),
		attribute.Int64("args.prod_slots", req.ProdSlots),
		attribute.String("args.sub_path", req.Subpath),
		attribute.String("args.prod_branch", req.ProdBranch),
		attribute.String("args.github_url", req.GithubUrl),
	)

	// Check the request is made by a user
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	// Find parent org
	org, err := s.admin.DB.FindOrganizationByName(ctx, req.OrganizationName)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check permissions
	if !claims.OrganizationPermissions(ctx, org.ID).CreateProjects {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to create projects")
	}

	// Check projects quota
	count, err := s.admin.DB.CountProjectsForOrganization(ctx, org.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if org.QuotaProjects >= 0 && count >= org.QuotaProjects {
		return nil, status.Errorf(codes.FailedPrecondition, "quota exceeded: org %q is limited to %d projects", org.Name, org.QuotaProjects)
	}

	// Check slots per deployment quota
	if org.QuotaSlotsPerDeployment >= 0 && int(req.ProdSlots) > org.QuotaSlotsPerDeployment {
		return nil, status.Errorf(codes.FailedPrecondition, "quota exceeded: org can't provision more than %d slots per deployment", org.QuotaSlotsPerDeployment)
	}

	// Check per project deployments and slots limit
	stats, err := s.admin.DB.CountDeploymentsForOrganization(ctx, org.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if org.QuotaDeployments >= 0 && stats.Deployments >= org.QuotaDeployments {
		return nil, status.Errorf(codes.FailedPrecondition, "quota exceeded: org %q is limited to %d deployments", org.Name, org.QuotaDeployments)
	}
	if org.QuotaSlotsTotal >= 0 && stats.Slots+int(req.ProdSlots) > org.QuotaSlotsTotal {
		return nil, status.Errorf(codes.FailedPrecondition, "quota exceeded: org %q is limited to %d total slots", org.Name, org.QuotaSlotsTotal)
	}

	// Check Github app is installed and caller has access on the repo
	installationID, err := s.getAndCheckGithubInstallationID(ctx, req.GithubUrl, claims.OwnerID())
	if err != nil {
		return nil, err
	}

	// Add prod TTL as 7 days if not a public project else infinite
	var prodTTL *int64
	if !req.Public {
		tmp := int64(prodDeplTTL.Seconds())
		prodTTL = &tmp
	}

	// Create the project
	proj, err := s.admin.CreateProject(ctx, org, claims.OwnerID(), &database.InsertProjectOptions{
		OrganizationID:       org.ID,
		Name:                 req.Name,
		Description:          req.Description,
		Public:               req.Public,
		Provisioner:          req.Provisioner,
		ProdVersion:          req.ProdVersion,
		ProdOLAPDriver:       req.ProdOlapDriver,
		ProdOLAPDSN:          req.ProdOlapDsn,
		ProdSlots:            int(req.ProdSlots),
		Subpath:              req.Subpath,
		ProdBranch:           req.ProdBranch,
		GithubURL:            &req.GithubUrl,
		GithubInstallationID: &installationID,
		ProdVariables:        req.Variables,
		ProdTTLSeconds:       prodTTL,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &adminv1.CreateProjectResponse{
		Project: s.projToDTO(proj, org.Name),
	}, nil
}

func (s *Server) DeleteProject(ctx context.Context, req *adminv1.DeleteProjectRequest) (*adminv1.DeleteProjectResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.OrganizationName),
		attribute.String("args.project", req.Name),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.OrganizationName, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProject {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to delete project")
	}

	err = s.admin.TeardownProject(ctx, proj)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &adminv1.DeleteProjectResponse{}, nil
}

func (s *Server) UpdateProject(ctx context.Context, req *adminv1.UpdateProjectRequest) (*adminv1.UpdateProjectResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.OrganizationName),
		attribute.String("args.project", req.Name),
	)
	if req.Description != nil {
		observability.AddRequestAttributes(ctx, attribute.String("args.description", *req.Description))
	}
	if req.Provisioner != nil {
		observability.AddRequestAttributes(ctx, attribute.String("args.provisioner", *req.Provisioner))
	}
	if req.ProdVersion != nil {
		observability.AddRequestAttributes(ctx, attribute.String("args.prod_version", *req.ProdVersion))
	}
	if req.ProdBranch != nil {
		observability.AddRequestAttributes(ctx, attribute.String("args.prod_branch", *req.ProdBranch))
	}
	if req.GithubUrl != nil {
		observability.AddRequestAttributes(ctx, attribute.String("args.github_url", *req.GithubUrl))
	}
	if req.Public != nil {
		observability.AddRequestAttributes(ctx, attribute.Bool("args.public", *req.Public))
	}
	if req.ProdSlots != nil {
		observability.AddRequestAttributes(ctx, attribute.Int64("args.prod_slots", *req.ProdSlots))
	}
	if req.ProdTtlSeconds != nil {
		observability.AddRequestAttributes(ctx, attribute.Int64("args.prod_ttl_seconds", *req.ProdTtlSeconds))
	}
	if req.NewName != nil {
		observability.AddRequestAttributes(ctx, attribute.String("args.new_name", *req.NewName))
	}

	// Check the request is made by a user
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	// Find project
	proj, err := s.admin.DB.FindProjectByName(ctx, req.OrganizationName, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProject {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to delete project")
	}

	githubURL := proj.GithubURL
	if req.GithubUrl != nil {
		// If changing the Github URL, check github app is installed and caller has access on the repo
		if safeStr(proj.GithubURL) != *req.GithubUrl {
			_, err = s.getAndCheckGithubInstallationID(ctx, *req.GithubUrl, claims.OwnerID())
			if err != nil {
				return nil, err
			}
			githubURL = req.GithubUrl
		}
	}

	prodTTLSeconds := proj.ProdTTLSeconds
	if req.ProdTtlSeconds != nil {
		if *req.ProdTtlSeconds == 0 {
			prodTTLSeconds = nil
		} else {
			prodTTLSeconds = req.ProdTtlSeconds
		}
	}

	opts := &database.UpdateProjectOptions{
		Name:                 valOrDefault(req.NewName, proj.Name),
		Description:          valOrDefault(req.Description, proj.Description),
		Public:               valOrDefault(req.Public, proj.Public),
		GithubURL:            githubURL,
		GithubInstallationID: proj.GithubInstallationID,
		ProdVersion:          valOrDefault(req.ProdVersion, proj.ProdVersion),
		ProdBranch:           valOrDefault(req.ProdBranch, proj.ProdBranch),
		ProdVariables:        proj.ProdVariables,
		ProdDeploymentID:     proj.ProdDeploymentID,
		ProdSlots:            int(valOrDefault(req.ProdSlots, int64(proj.ProdSlots))),
		ProdTTLSeconds:       prodTTLSeconds,
		Provisioner:          valOrDefault(req.Provisioner, proj.Provisioner),
		Annotations:          proj.Annotations,
	}
	proj, err = s.admin.UpdateProject(ctx, proj, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.UpdateProjectResponse{
		Project: s.projToDTO(proj, req.OrganizationName),
	}, nil
}

func (s *Server) GetProjectVariables(ctx context.Context, req *adminv1.GetProjectVariablesRequest) (*adminv1.GetProjectVariablesResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.OrganizationName),
		attribute.String("args.project", req.Name),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.OrganizationName, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProject {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project variables")
	}

	return &adminv1.GetProjectVariablesResponse{Variables: proj.ProdVariables}, nil
}

func (s *Server) UpdateProjectVariables(ctx context.Context, req *adminv1.UpdateProjectVariablesRequest) (*adminv1.UpdateProjectVariablesResponse, error) {
	proj, err := s.admin.DB.FindProjectByName(ctx, req.OrganizationName, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check the request is made by a user
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProject {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to update project variables")
	}

	proj, err = s.admin.UpdateProject(ctx, proj, &database.UpdateProjectOptions{
		Name:                 proj.Name,
		Description:          proj.Description,
		Public:               proj.Public,
		GithubURL:            proj.GithubURL,
		GithubInstallationID: proj.GithubInstallationID,
		ProdVersion:          proj.ProdVersion,
		ProdBranch:           proj.ProdBranch,
		ProdVariables:        req.Variables,
		ProdDeploymentID:     proj.ProdDeploymentID,
		ProdSlots:            proj.ProdSlots,
		ProdTTLSeconds:       proj.ProdTTLSeconds,
		Provisioner:          proj.Provisioner,
		Annotations:          proj.Annotations,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "variables updated failed with error %s", err.Error())
	}

	return &adminv1.UpdateProjectVariablesResponse{Variables: proj.ProdVariables}, nil
}

func (s *Server) ListProjectMembers(ctx context.Context, req *adminv1.ListProjectMembersRequest) (*adminv1.ListProjectMembersResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.project", req.Project),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.Superuser(ctx) && !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ReadProjectMembers {
		return nil, status.Error(codes.PermissionDenied, "not authorized to read project members")
	}

	token, err := unmarshalPageToken(req.PageToken)
	if err != nil {
		return nil, err
	}
	pageSize := validPageSize(req.PageSize)

	members, err := s.admin.DB.FindProjectMemberUsers(ctx, proj.ID, token.Val, pageSize)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	nextToken := ""
	if len(members) >= pageSize {
		nextToken = marshalPageToken(members[len(members)-1].Email)
	}

	dtos := make([]*adminv1.Member, len(members))
	for i, member := range members {
		dtos[i] = memberToPB(member)
	}

	return &adminv1.ListProjectMembersResponse{
		Members:       dtos,
		NextPageToken: nextToken,
	}, nil
}

func (s *Server) ListProjectInvites(ctx context.Context, req *adminv1.ListProjectInvitesRequest) (*adminv1.ListProjectInvitesResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.project", req.Project),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ReadProjectMembers {
		return nil, status.Error(codes.PermissionDenied, "not authorized to read project members")
	}

	token, err := unmarshalPageToken(req.PageToken)
	if err != nil {
		return nil, err
	}
	pageSize := validPageSize(req.PageSize)

	// get pending user invites for this project
	userInvites, err := s.admin.DB.FindProjectInvites(ctx, proj.ID, token.Val, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	nextToken := ""
	if len(userInvites) >= pageSize {
		nextToken = marshalPageToken(userInvites[len(userInvites)-1].Email)
	}

	invitesDtos := make([]*adminv1.UserInvite, len(userInvites))
	for i, invite := range userInvites {
		invitesDtos[i] = inviteToPB(invite)
	}

	return &adminv1.ListProjectInvitesResponse{
		Invites:       invitesDtos,
		NextPageToken: nextToken,
	}, nil
}

func (s *Server) AddProjectMember(ctx context.Context, req *adminv1.AddProjectMemberRequest) (*adminv1.AddProjectMemberResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.String("args.role", req.Role),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProjectMembers {
		return nil, status.Error(codes.PermissionDenied, "not allowed to add project members")
	}

	// Check outstanding invites quota
	count, err := s.admin.DB.CountInvitesForOrganization(ctx, proj.OrganizationID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	org, err := s.admin.DB.FindOrganization(ctx, proj.OrganizationID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if org.QuotaOutstandingInvites >= 0 && count >= org.QuotaOutstandingInvites {
		return nil, status.Errorf(codes.FailedPrecondition, "quota exceeded: org %q can at most have %d outstanding invitations", org.Name, org.QuotaOutstandingInvites)
	}

	role, err := s.admin.DB.FindProjectRole(ctx, req.Role)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var invitedByUserID, invitedByName string
	if claims.OwnerType() == auth.OwnerTypeUser {
		user, err := s.admin.DB.FindUser(ctx, claims.OwnerID())
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		invitedByUserID = user.ID
		invitedByName = user.DisplayName
	}

	user, err := s.admin.DB.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Invite user to join the project
		err := s.admin.DB.InsertProjectInvite(ctx, &database.InsertProjectInviteOptions{
			Email:     req.Email,
			InviterID: invitedByUserID,
			ProjectID: proj.ID,
			RoleID:    role.ID,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Send invitation email
		err = s.admin.Email.SendProjectInvite(&email.ProjectInvite{
			ToEmail:       req.Email,
			ToName:        "",
			AdminURL:      s.opts.ExternalURL,
			FrontendURL:   s.opts.FrontendURL,
			OrgName:       org.Name,
			ProjectName:   proj.Name,
			RoleName:      role.Name,
			InvitedByName: invitedByName,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &adminv1.AddProjectMemberResponse{
			PendingSignup: true,
		}, nil
	}

	err = s.admin.DB.InsertProjectMemberUser(ctx, proj.ID, user.ID, role.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.admin.Email.SendProjectAddition(&email.ProjectAddition{
		ToEmail:       req.Email,
		ToName:        "",
		FrontendURL:   s.opts.FrontendURL,
		OrgName:       org.Name,
		ProjectName:   proj.Name,
		RoleName:      role.Name,
		InvitedByName: invitedByName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.AddProjectMemberResponse{
		PendingSignup: false,
	}, nil
}

func (s *Server) RemoveProjectMember(ctx context.Context, req *adminv1.RemoveProjectMemberRequest) (*adminv1.RemoveProjectMemberResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.project", req.Project),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProjectMembers {
		return nil, status.Error(codes.PermissionDenied, "not allowed to remove project members")
	}

	user, err := s.admin.DB.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		// check if there is a pending invite
		invite, err := s.admin.DB.FindProjectInvite(ctx, proj.ID, req.Email)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return nil, status.Error(codes.InvalidArgument, "user not found")
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		err = s.admin.DB.DeleteProjectInvite(ctx, invite.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &adminv1.RemoveProjectMemberResponse{}, nil
	}

	err = s.admin.DB.DeleteProjectMemberUser(ctx, proj.ID, user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.RemoveProjectMemberResponse{}, nil
}

func (s *Server) SetProjectMemberRole(ctx context.Context, req *adminv1.SetProjectMemberRoleRequest) (*adminv1.SetProjectMemberRoleResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.String("args.role", req.Role),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ManageProjectMembers {
		return nil, status.Error(codes.PermissionDenied, "not allowed to set project member roles")
	}

	role, err := s.admin.DB.FindProjectRole(ctx, req.Role)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.admin.DB.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		// Check if there is a pending invite for this user
		invite, err := s.admin.DB.FindProjectInvite(ctx, proj.ID, req.Email)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return nil, status.Error(codes.InvalidArgument, "user not found")
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		err = s.admin.DB.UpdateProjectInviteRole(ctx, invite.ID, role.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &adminv1.SetProjectMemberRoleResponse{}, nil
	}

	err = s.admin.DB.UpdateProjectMemberUserRole(ctx, proj.ID, user.ID, role.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.SetProjectMemberRoleResponse{}, nil
}

// getAndCheckGithubInstallationID returns a valid installation ID iff app is installed and user is a collaborator of the repo
func (s *Server) getAndCheckGithubInstallationID(ctx context.Context, githubURL, userID string) (int64, error) {
	// Get Github installation ID for the repo
	installationID, err := s.admin.GetGithubInstallation(ctx, githubURL)
	if err != nil {
		if errors.Is(err, admin.ErrGithubInstallationNotFound) {
			return 0, status.Errorf(codes.PermissionDenied, "you have not granted Rill access to %q", githubURL)
		}

		return 0, status.Errorf(codes.Internal, "failed to get Github installation: %q", err.Error())
	}

	if installationID == 0 {
		return 0, status.Errorf(codes.PermissionDenied, "you have not granted Rill access to %q", githubURL)
	}

	// Check that user is a collaborator on the repo
	user, err := s.admin.DB.FindUser(ctx, userID)
	if err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}

	if user.GithubUsername == "" {
		return 0, status.Errorf(codes.PermissionDenied, "you have not granted Rill access to your Github account")
	}

	_, err = s.admin.LookupGithubRepoForUser(ctx, installationID, githubURL, user.GithubUsername)
	if err != nil {
		if errors.Is(err, admin.ErrUserIsNotCollaborator) {
			return 0, status.Errorf(codes.PermissionDenied, "you are not collaborator to the repo %q", githubURL)
		}
		return 0, status.Error(codes.Internal, err.Error())
	}

	return installationID, nil
}

// SudoUpdateTags updates the tags for a project in organization for superusers
func (s *Server) SudoUpdateAnnotations(ctx context.Context, req *adminv1.SudoUpdateAnnotationsRequest) (*adminv1.SudoUpdateAnnotationsResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.Int("args.annotations", len(req.Annotations)),
	)

	// Check the request is made by a superuser
	claims := auth.GetClaims(ctx)
	if !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to update annotations")
	}

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	proj, err = s.admin.UpdateProject(ctx, proj, &database.UpdateProjectOptions{
		Name:                 proj.Name,
		Description:          proj.Description,
		Public:               proj.Public,
		GithubURL:            proj.GithubURL,
		GithubInstallationID: proj.GithubInstallationID,
		ProdVersion:          proj.ProdVersion,
		ProdBranch:           proj.ProdBranch,
		ProdVariables:        proj.ProdVariables,
		ProdDeploymentID:     proj.ProdDeploymentID,
		ProdSlots:            proj.ProdSlots,
		ProdTTLSeconds:       proj.ProdTTLSeconds,
		Provisioner:          proj.Provisioner,
		Annotations:          req.Annotations,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.SudoUpdateAnnotationsResponse{
		Project: s.projToDTO(proj, req.Organization),
	}, nil
}

func (s *Server) projToDTO(p *database.Project, orgName string) *adminv1.Project {
	frontendURL, _ := url.JoinPath(s.opts.FrontendURL, orgName, p.Name)

	return &adminv1.Project{
		Id:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		Public:           p.Public,
		OrgId:            p.OrganizationID,
		OrgName:          orgName,
		Provisioner:      p.Provisioner,
		ProdVersion:      p.ProdVersion,
		ProdOlapDriver:   p.ProdOLAPDriver,
		ProdOlapDsn:      p.ProdOLAPDSN,
		ProdSlots:        int64(p.ProdSlots),
		ProdBranch:       p.ProdBranch,
		Subpath:          p.Subpath,
		GithubUrl:        safeStr(p.GithubURL),
		ProdDeploymentId: safeStr(p.ProdDeploymentID),
		ProdTtlSeconds:   safeInt64(p.ProdTTLSeconds),
		FrontendUrl:      frontendURL,
		Annotations:      p.Annotations,
		CreatedOn:        timestamppb.New(p.CreatedOn),
		UpdatedOn:        timestamppb.New(p.UpdatedOn),
	}
}

func deploymentToDTO(d *database.Deployment) *adminv1.Deployment {
	var s adminv1.DeploymentStatus
	switch d.Status {
	case database.DeploymentStatusUnspecified:
		s = adminv1.DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED
	case database.DeploymentStatusPending:
		s = adminv1.DeploymentStatus_DEPLOYMENT_STATUS_PENDING
	case database.DeploymentStatusOK:
		s = adminv1.DeploymentStatus_DEPLOYMENT_STATUS_OK
	case database.DeploymentStatusError:
		s = adminv1.DeploymentStatus_DEPLOYMENT_STATUS_ERROR
	default:
		panic(fmt.Errorf("unhandled deployment status %d", d.Status))
	}

	return &adminv1.Deployment{
		Id:                d.ID,
		ProjectId:         d.ProjectID,
		Slots:             int64(d.Slots),
		Branch:            d.Branch,
		RuntimeHost:       d.RuntimeHost,
		RuntimeInstanceId: d.RuntimeInstanceID,
		Status:            s,
		StatusMessage:     d.StatusMessage,
		CreatedOn:         timestamppb.New(d.CreatedOn),
		UpdatedOn:         timestamppb.New(d.UpdatedOn),
	}
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeInt64(s *int64) int64 {
	if s == nil {
		return 0
	}
	return *s
}

func valOrDefault[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
