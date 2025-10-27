package ddqp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GenericParser_Parse(t *testing.T) {
	parser := NewGenericParser()

	tests := []struct {
		name     string
		query    string
		wantErr  bool
		printAST bool // For debugging, can opt in to print AST
	}{
		{
			name:     "addition",
			query:    "sum:metric.name{foo:bar} + sum:metric.name_two{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "subtraction",
			query:    "sum:metric.name{foo:bar} - sum:metric.name_two{foo:bar} - 0.1",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "multiplication",
			query:    "sum:metric.name{foo:bar} * sum:metric.name_two{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "division",
			query:    "sum:metric.name{foo:bar} / sum:metric.name_two{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "parens",
			query:    "(sum:metric.name{foo:bar} - sum:metric.name_two{foo:bar}) / sum:metric.name_two{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "filters with slashes",
			query:    "sum:metric.name{foo:bar/hello} / sum:metric.name_two{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "parens division with int multiplication",
			query:    "(sum:metric.name{foo:bar/hello} / sum:metric.name_two{baz:bang}) / 100",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "complex expression with multiple operators",
			query:    "sum:metric.name{foo:bar} * sum:metric.name_two{foo:bar} + sum:metric.name_three{foo:bar} - 50",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "nested parentheses",
			query:    "(sum:metric.name{foo:bar} - (sum:metric.name_two{foo:bar} * 2)) / 10",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "complex filters in expression",
			query:    "sum:metric.name{env:prod AND service:api} - sum:metric.name{env:staging AND service:api}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "multiple arithmetic operations with constants",
			query:    "sum:metric.name{*} * 2",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "expression with grouping in metrics",
			query:    "sum:metric.name{foo:bar} / sum:metric.name_two{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "expression with numeric values in filters",
			query:    "sum:metric.name{code:200} / sum:metric.name_two{code:200}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "simple query",
			query:    "moving_rollup(default_zero(sum:metric.name{app:bazz,env:staging}.as_rate()), 60, 'avg')",
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
			name:     "test complex query with wildcard filter",
			query:    "metric.name{app:bazz,env:staging,host:host-*}.as_rate()",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "wrapping function around metric expression",
			query:    "default_zero(avg:metric.name{foo:bar} + avg:other.metric.name{foo:bar})",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "wrapping function around one metric inside of the metric expression",
			query:    "default_zero(avg:metric.name{foo:bar}) + avg:other.metric.name{foo:bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test a really weird complex query",
			query:    "default_zero(avg:system.cpu.user{foo:bar, (kube_cluster_name:test-cluster OR !kube_cluster_name:*)}.rollup(avg, 300)) + (avg:system.cpu.user{foo:bar, (kube_cluster_name:test-cluster OR !kube_cluster_name:*)} * 1000 + (100 / 10))",
			wantErr:  false,
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
				want, got := sanitizeWantAndGot(tt.query, ast.String())
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_GenericParser_FromFile(t *testing.T) {
	parser := NewGenericParser()

	f, err := os.Open("./test_queries.txt")
	if err != nil {
		t.Skipf("skipping file-driven tests; could not open test_queries.txt: %v", err)
		return
	}
	t.Cleanup(func() {
		require.NoError(t, f.Close())
	})

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		t.Run("line_"+fmt.Sprintf("%d", lineNum), func(t *testing.T) {
			ast, err := parser.Parse(line)
			require.NoError(t, err, "failed to parse line %d: %s", lineNum, line)
			want, got := sanitizeWantAndGot(line, ast.String())
			assert.Equal(t, want, got)
		})
	}

	require.NoError(t, scanner.Err())
}

// we do not care about spacing in most places, but do want to make sure that in cases like AND NOT that
// the query is actually parsing correctly and being returned with the correct spacing.
func sanitizeWantAndGot(want, got string) (string, string) {
	want = strings.ReplaceAll(want, ", ", ",")
	got = strings.ReplaceAll(got, ", ", ",")
	want = strings.ReplaceAll(want, " }", "}")
	got = strings.ReplaceAll(got, " }", "}")
	want = strings.ReplaceAll(want, "{ ", "{")
	got = strings.ReplaceAll(got, "{ ", "{")
	want = strings.ToLower(want)
	got = strings.ToLower(got)
	return want, got
}
