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
