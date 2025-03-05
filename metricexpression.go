package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Operator int

const (
	OpMul Operator = iota
	OpDiv
	OpAdd
	OpSub
)

var operatorMap = map[string]Operator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv}

func (o *Operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

type ExprValue struct {
	Subexpression *MetricExpression `  "(" @@ ")"`
	MetricQuery   *MetricQuery      `| @@`
	Number        *float64          `| @Ident`
}

func (expr *ExprValue) GetQueries() []string {
	strs := []string{}
	if expr.Subexpression != nil {
		m := expr.Subexpression.GetQueries()
		for _, v := range m {
			strs = append(strs, v)
		}
		return strs
	}

	if expr.MetricQuery != nil {
		strs = append(strs, expr.MetricQuery.String())
		return strs
	}

	return []string{}
}

type Factor struct {
	Base *ExprValue `@@`
}

func (f *Factor) GetQueries() []string {
	return f.Base.GetQueries()
}

type OpFactor struct {
	Operator Operator `@("*" | "/")`
	Factor   *Factor  `@@`
}

type Term struct {
	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

func (t *Term) GetQueries() []string {
	queries := t.Left.GetQueries()

	for _, v := range t.Right {
		queries = append(queries, v.Factor.GetQueries()...)
	}

	return queries
}

type OpTerm struct {
	Operator Operator `@("+" | "-")`
	Term     *Term    `@@`
}

type MetricExpression struct {
	Pos lexer.Position

	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

func (me *MetricExpression) GetQueries() map[string]string {
	queries := me.Left.GetQueries()

	for _, v := range me.Right {
		rightQueries := v.Term.GetQueries()
		queries = append(queries, rightQueries...)
	}

	// Create a map to store unique queries in order of appearance
	queryMap := make(map[string]string)
	seen := make(map[string]bool)
	orderedQueries := make([]string, 0)

	// First pass: collect unique queries in order
	for _, query := range queries {
		if !seen[query] {
			seen[query] = true
			orderedQueries = append(orderedQueries, query)
		}
	}

	// Second pass: assign variables in order
	for i, query := range orderedQueries {
		queryMap[toCharStr(i+1)] = query
	}

	return queryMap
}

func toCharStr(i int) string {
	const abc = "abcdefghijklmnopqrstuvwxyz"
	return abc[i-1 : i]
}

// NewMetricExpressionParser returns a Parser which is capable of interpretting
// a metric expression.
func NewMetricExpressionParser() *MetricExpressionParser {
	mep := &MetricExpressionParser{
		parser: participle.MustBuild[MetricExpression](
			participle.Lexer(lex),
			participle.Unquote("String"),
		),
	}

	return mep
}

// Display

func (o Operator) String() string {
	switch o {
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpSub:
		return "-"
	case OpAdd:
		return "+"
	}
	panic("unsupported operator")
}

func (expr *ExprValue) String() string {
	if expr.Number != nil {
		return fmt.Sprintf("%g", *expr.Number)
	}
	if expr.MetricQuery != nil {
		return expr.MetricQuery.String()
	}
	return "(" + expr.Subexpression.String() + ")"
}

func (f *Factor) String() string {
	out := f.Base.String()
	return out
}

func (o *OpFactor) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Factor)
}

func (t *Term) String() string {
	out := []string{t.Left.String()}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

func (me *MetricExpression) String() string {
	out := []string{me.Left.String()}
	for _, r := range me.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

// MetricExpressionParser is parser returned when calling NewMetricExpressionParser.
type MetricExpressionParser struct {
	parser *participle.Parser[MetricExpression]
}

// Parse sanitizes the query string and returns the AST and any error.
func (mep *MetricExpressionParser) Parse(expr string) (*MetricExpression, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(expr, "\n", "")
	// return the raw parsed outpu
	return mep.parser.ParseString("", sanitized)
}

// MetricExpressionFormula breaks down a query into its formulaic parts
type MetricExpressionFormula struct {
	Expressions map[string]string
	Formula     string
}

func NewMetricExpressionFormula(expr *MetricExpression) *MetricExpressionFormula {
	exprFormula := &MetricExpressionFormula{
		Expressions: map[string]string{},
		Formula:     "",
	}

	exprFormula.Expressions = expr.GetQueries()

	formula := expr.String()

	for key, value := range exprFormula.Expressions {
		formula = strings.ReplaceAll(formula, value, key)
	}

	exprFormula.Formula = formula
	return exprFormula
}
