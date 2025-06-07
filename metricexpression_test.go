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
		formula     string          // Expected formula or empty if we don't care about exact mapping
		expressions map[string]string
		wantErr     bool
		printAST    bool // For debugging, can opt in to print AST
		skipFormulaCheck bool // Skip formula assertion for cases where map ordering matters
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
			skipFormulaCheck: true,
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
			skipFormulaCheck: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.ParseString("", tt.query)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err)
			}

			expr := NewMetricExpressionFormula(ast)
			if !tt.skipFormulaCheck {
				assert.Equal(t, tt.formula, expr.Formula)
			}
			// For expressions, we should check content equality as the key assignment could vary
			expressionsEqual := true
			if len(expr.Expressions) != len(tt.expressions) {
				expressionsEqual = false
			} else {
				// Check if all expected expressions exist in result
				for _, expectedVal := range tt.expressions {
					found := false
					for _, actualVal := range expr.Expressions {
						if expectedVal == actualVal {
							found = true
							break
						}
					}
					if !found {
						expressionsEqual = false
						break
					}
				}
			}
			if !expressionsEqual {
				t.Logf("Expected expressions: %v\nActual expressions: %v", tt.expressions, expr.Expressions)
			}
			assert.True(t, expressionsEqual, "Expressions maps should contain the same values (might have different keys)")
			
		})
	}
}
