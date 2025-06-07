package ddqp

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMetricFilterParser() *participle.Parser[MetricFilter] {
	return participle.MustBuild[MetricFilter](
		participle.Lexer(lex),
		participle.Unquote("String"),
	)
}

func Test_MetricMonitorFilter(t *testing.T) {
	parser := newMetricFilterParser()

	tests := []struct {
		name     string
		query    string
		wantErr  bool
		printAST bool // For debugging, can opt in to print AST
	}{
		{
			name:     "test asterisk only",
			query:    "*",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test int and string",
			query:    "code:2xx",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test one simple filter",
			query:    "foo:bar-bar",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple comma separated filter",
			query:    "a:b, c:d",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple AND separated filter",
			query:    "a:b AND c:d AND e:f",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple AND NOT separated filter",
			query:    "a:b AND c:d AND NOT e:f",
			wantErr:  false,
			printAST: true,
		},
		{
			name:     "test simple OR separated filter",
			query:    "a:b OR c:d OR e:f",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple parens filter",
			query:    "c IN (d)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple parens filter with comma separated values",
			query:    "a IN (b, c, d)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple parens filter with OR separated values",
			query:    "e IN (f OR g OR h)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test negative filter with !",
			query:    "!a:b",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test multiple negative filter with !",
			query:    "!a:b, !c:d",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test AND then parens",
			query:    "a:b AND (c:d OR e:f)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test example from DataDog docs",
			query:    "env:shop.ist AND availability-zone IN (us-east-1a, us-east-1b, us-east4-b)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test another example from DataDog docs",
			query:    "env:prod AND location NOT IN (atlanta, seattle, las-vegas)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test complex nested logical operators",
			query:    "service:api AND ((env:prod AND region:us-east) OR (env:staging AND region:us-west))",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test filter with special characters",
			query:    "path:api_endpoint, method:GET",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test filter with numeric values",
			query:    "status:200, response_time:500",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test multiple IN operators combined",
			query:    "env IN (prod, staging) AND service IN (web, api, worker)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test NOT operator with complex expression",
			query:    "env:prod AND NOT (region:us-east AND datacenter:primary)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test wildcard in filter values",
			query:    "host:web-*, service:api-*-service",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test deeply nested expressions with mixed operators",
			query:    "(service:api AND (env:prod OR env:staging)) OR (service:web AND env:dev AND NOT region:eu-west)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical greater than comparison",
			query:    "response_time:>500",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical less than comparison",
			query:    "error_rate:<0.01",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical greater than or equal comparison",
			query:    "cpu:>=90",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical less than or equal comparison",
			query:    "memory:<=75.5",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical comparison with complex expression",
			query:    "response_time:>500 AND status:200",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test multiple numerical comparisons",
			query:    "cpu:>80 OR memory:>90",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical comparison in nested expression",
			query:    "env:prod AND (response_time:<200 OR error_rate:<0.01)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test greater than comparison",
			query:    "response_time:>500",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test less than comparison",
			query:    "error_rate:<0.01",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test greater than or equal comparison",
			query:    "cpu:>=90",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test less than or equal comparison",
			query:    "memory:<=75.5",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical comparison in complex expressions",
			query:    "response_time:>500 AND status:200",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test multiple numerical comparisons",
			query:    "cpu:>80 OR memory:>90",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numerical comparison in nested expression",
			query:    "env:prod AND (response_time:<200 OR error_rate:<0.01)",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test regex filter simple",
			query:    "service:~simple-regex",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test regex filter with API version pattern",
			query:    "path:~simple-pattern",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test regex filter with AND operator",
			query:    "service:~api-.* AND env:prod",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test regex filter with OR operator",
			query:    "service:~simple-regex OR env:~simple-pattern",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test regex filter in nested expression",
			query:    "env:prod AND (service:~api-.* OR host:~web-.*)",
			wantErr:  false,
			printAST: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.ParseString("", tt.query)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err)
			}

			if tt.printAST {
				repr.Println(ast)
			}

			// Check to make sure we're able to restringify
			assert.Equal(t, tt.query, ast.String())
		})
	}
}