package resolvers

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/rilldata/rill/runtime"
	metricssqlparser "github.com/rilldata/rill/runtime/pkg/metricssql"
)

func init() {
	runtime.RegisterResolverInitializer("metrics_sql", newMetricsSQL)
}

type metricsSQLProps struct {
	SQL string `mapstructure:"sql"`
}

// newMetricsSQL creates a resolver for evaluating metrics SQL.
// It wraps the regular SQL resolver and compiles the metrics SQL to a regular SQL query first.
// The compiler preserves templating in the SQL, allowing the regular SQL resolver to handle SQL templating rules.
func newMetricsSQL(ctx context.Context, opts *runtime.ResolverOptions) (runtime.Resolver, error) {
	props := &metricsSQLProps{}
	if err := mapstructure.Decode(opts.Properties, props); err != nil {
		return nil, err
	}

	if props.SQL == "" {
		return nil, errors.New(`metrics SQL: missing required property "sql"`)
	}

	ctrl, err := opts.Runtime.Controller(ctx, opts.InstanceID)
	if err != nil {
		return nil, err
	}

	compiler := metricssqlparser.New(ctrl, opts.InstanceID, opts.UserAttributes)
	sql, connector, refs, err := compiler.Compile(ctx, props.SQL)
	if err != nil {
		return nil, err
	}

	// Build the options for the regular SQL resolver
	sqlResolverOpts := &runtime.ResolverOptions{
		Runtime:    opts.Runtime,
		InstanceID: opts.InstanceID,
		Properties: map[string]any{
			"connector": connector,
			"sql":       sql,
		},
		Args:           opts.Args,
		UserAttributes: opts.UserAttributes,
		ForExport:      opts.ForExport,
	}

	return newSQLWithRefs(ctx, sqlResolverOpts, refs)
}
