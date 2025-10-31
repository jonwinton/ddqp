package ddqp

import (
	"strings"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MetricQuery(t *testing.T) {
	parser := NewMetricQueryParser()

	tests := []struct {
		name     string
		query    string
		wantErr  bool
		printAST bool // For debugging, can opt in to print AST
	}{
		{
			name:     "simple query",
			query:    "moving_rollup(default_zero(sum:metric{key:value,!service:service,env:staging}.as_rate()), 60, 'avg')",
			wantErr:  false,
			printAST: false,
		},
		// Simple passing example. Guaranteed to have all parts of a query
		{
			name:     "simple query",
			query:    "sum:namespace.metric.name{foo:bar} by {foo}",
			wantErr:  false,
			printAST: false,
		},
		// Simple failing example. Guaranteed to fail because missing the aggregator
		{
			name:     "fail due to no aggregator",
			query:    "namespace.metric.name{foo:bar, baz:bang} by {foo,bar}",
			wantErr:  true,
			printAST: false,
		},
		{
			name:     "filter by asterisk",
			query:    "sum:namespace.metric.name{*} by {foo,bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "filer by partial asterisk",
			query:    "sum:namespace.metric.name{foo:bar-*} by {foo,bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test underscores in metric name",
			query:    "sum:namespace.metric_name{foo:bar} by {baz}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test hyphens in filters and groupings",
			query:    "sum:prometheus_metric_source{foo:bar-bar, baz:bang} by {fizz-buzz,bang}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numbers in the metric name",
			query:    "sum:prometheus_metric_source_1{foo:bar-bar, baz:bang} by {fizz-buzz,bang}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numbers in the filters and grouping",
			query:    "sum:prometheus_metric_source_1{foo:bar-bar-1, baz:bang_2} by {fizz-buzz3,bang}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test docs example query",
			query:    "avg:system.cpu.user{env:staging AND (availability-zone:us-east-1a OR availability-zone:us-east-1c)} by {availability-zone}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test less than condition in count",
			query:    "count(v: v<=1):metric.name{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test greater than condition in count",
			query:    "count(v: v>=1.53):metric.name{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test equal than condition in count",
			query:    "count(v: v>=100):metric.name{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "function with no args",
			query:    "sum:system.cpu.user{*}.as_rate()",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "function with identifier and number args",
			query:    "sum:system.cpu.user{*}.rollup(avg,60)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "function with string arg",
			query:    "sum:system.cpu.user{*}.label(\"CPU User\")",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "function with boolean arg",
			query:    "sum:system.cpu.user{*}.fill(true)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "chained functions mixed args",
			query:    "sum:system.cpu.user{*}.as_rate().rollup(avg,60).label(\"CPU User\").fill(true)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "count with less than condition",
			query:    "count(v: v<10):metric.name{*}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "negated simple filter",
			query:    "sum:metric.name{!env:prod, region:us-east-1}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "regex filter with string",
			query:    "sum:metric.name{host:~\"web-.*\"}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "comparison filters with AND/OR",
			query:    "sum:metric.name{duration:>=100 AND duration:<=200 OR errors:>5}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "IN list filter",
			query:    "sum:metric.name{env IN (prod, staging)}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "NOT IN list filter",
			query:    "sum:metric.name{region NOT IN (us-east-1, us-west-2)}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "grouped filter with AND NOT and OR",
			query:    "sum:metric.name{(service:api AND NOT env:dev) OR region IN (us-east-1, us-west-2)}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "wildcard and slash in metric name with functions and grouping",
			query:    "sum:system.disk/*{*} by {host}.as_rate().rollup(avg,300)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "query that starts with a metric name",
			query:    "metric{filter:value-*}.as_rate()",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "wildcard filter",
			query:    "sum:metric{key:*}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "filter starts with wildcard and ends with wildcard",
			query:    "avg:metric{key:*value-*}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "wildcard by",
			query:    "avg:metric{key:*value-*} by {*}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "wildcard by error",
			query:    "avg:metric{key:*value-*} by {*foo}",
			wantErr:  true,
			printAST: false,
		},
		{
			name:     "no filters",
			query:    "sum:requests.count",
			wantErr:  true,
			printAST: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.Parse(tt.query)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err)
			}

			if tt.printAST {
				repr.Println(ast)
			}

			// Check to make sure we're able to restringify
			if !tt.wantErr {
				want := strings.ReplaceAll(tt.query, ", ", ",")
				got := strings.ReplaceAll(ast.String(), ", ", ",")
				assert.Equal(t, want, got)
			}
		})
	}
}
