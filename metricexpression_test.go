package ddqp

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMetricExpressionParser() *participle.Parser[MetricExpression] {
	return participle.MustBuild[MetricExpression](
		participle.Lexer(lex),
		participle.Unquote("String"),
	)
}

func Test_MetricExpressionCanParse(t *testing.T) {
	parser := newMetricExpressionParser()

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

			// Assert equal stingification
			assert.Equal(t, tt.query, ast.String())
		})
	}
}

func Test_MetricExpressionFormula(t *testing.T) {
	parser := newMetricExpressionParser()

	tests := []struct {
		name        string
		query       string
		formula     string
		expressions map[string]string
		wantErr     bool
		printAST    bool // For debugging, can opt in to print AST
	}{
		{
			name:    "addition formula",
			query:   "sum:metric.name{foo:bar} + sum:metric.name_two{foo:bar, baz:bang}",
			formula: "a + b",
			expressions: map[string]string{
				"a": "sum:metric.name{foo:bar}",
				"b": "sum:metric.name_two{foo:bar, baz:bang}",
			},
			wantErr:  false,
			printAST: false,
		},
		{
			name:    "addition and division formula",
			query:   "(sum:metric.name{foo:bar} + sum:metric.name_two{foo:bar}) / sum:metric.name_three{*}",
			formula: "(a + b) / c",
			expressions: map[string]string{
				"a": "sum:metric.name{foo:bar}",
				"b": "sum:metric.name_two{foo:bar}",
				"c": "sum:metric.name_three{*}",
			},
			wantErr:  false,
			printAST: false,
		},
		{
			name:    "calculate percent",
			query:   "(sum:metric.name{foo:bar} / sum:metric.name_two{foo:bar}) * 100",
			formula: "(a / b) * 100",
			expressions: map[string]string{
				"a": "sum:metric.name{foo:bar}",
				"b": "sum:metric.name_two{foo:bar}",
			},
			wantErr:  false,
			printAST: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.ParseString("", tt.query)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err)
			}

			expr := NewMetricExpressionFormula(ast)
			assert.Equal(t, tt.formula, expr.Formula)
			assert.Equal(t, tt.expressions, expr.Expressions)
		})
	}
}
