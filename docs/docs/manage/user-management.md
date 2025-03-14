---
title: User Management
sidebar_label: User Management
sidebar_position: 20
---

In Rill Cloud, access can be granted at the organization or project level using the Rill CLI.

## Install and authenticate the Rill CLI

To manage cloud permissions with the Rill CLI, you must first authenticate it. If you have not already done so, run:
```
rill login
```

## Managing members of an organization

When you invite a user to an organization on Rill Cloud, they automatically get access to *all projects* in the organization. Users can have one of two roles:

- **Viewers** can browse projects and view dashboards
- **Admins** can manage projects by deploying new projects, making changes to existing projects, or deleting deployed projects. They can also manage members of an organization by granting or revoking access to other users.  
  
### Add a member

To add a member to an organization, run the following command:
```
rill user add
```
You will then be prompted for details about the user.

If you add a user who has not yet signed up for Rill, they will receive an email inviting them to sign up and join.

### Automatically add members by email domain

You can automatically add users to your organization by their email domain. For example, if you whitelist `yourdomain.com`, new and existing users with an email address ending on `@yourdomain.com` will automatically be added to your organization.

The feature currently requires manual action by a support representative at Rill. Just [reach out here](https://www.rilldata.com/contact) and ask us to whitelist your domain.

### Other actions

Run `rill user --help` to show commands for listing members or changing access.

## Managing members of a project

By default, adding a user to an organization grants them access to all its projects. You can alternatively add a user only to a specific project. Users can have one of two roles on a project:

- **Viewers** can view the project's dashboards
- **Admins** can additionally edit the project, and view and edit project members

### Add a member

To add a member to a project, run the following command:
```
rill user add --project [PROJECT NAME]
```
You will then be prompted for details about the user. HINT: Run `rill project list` to show available projects.

If you add a user who has not yet signed up for Rill, they will receive an email inviting them to join.

### Other actions

Run `rill user --help` to show commands for listing members or changing access.

## Make a project public

Projects on Rill Cloud are private by default. To make a project's dashboards publicly accessible without authentication, run:
```
rill project edit --public=true
```

:::caution Avoid Sharing Private Data

Warning: If you make a project public, make sure it does not expose any confidential data.

:::
