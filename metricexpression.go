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
	Number        *float64          `  @(Float|Int)`
	MetricQuery   *MetricQuery      `| @@`
	Subexpression *MetricExpression `| "(" @@ ")"`
}

type Factor struct {
	Base *ExprValue `@@`
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
	queries := []string{t.Left.String()}

	for _, v := range t.Right {
		queries = append(queries, v.Factor.String())
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

	queryMap := make(map[string]string)
	for key, value := range queries {
		queryMap[fmt.Sprintf("%d", key)] = value
		// queryMap[toCharStr(key+1)] = value
	}
	return queryMap
}

// func toCharStr(i int) string {
// 	return strings.ToLower(string('A' - 1 + i))
// }

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

func (v *ExprValue) String() string {
	if v.Number != nil {
		return fmt.Sprintf("%g", *v.Number)
	}
	if v.MetricQuery != nil {
		return v.MetricQuery.String()
	}
	return "(" + v.Subexpression.String() + ")"
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

func (e *MetricExpression) String() string {
	out := []string{e.Left.String()}
	for _, r := range e.Right {
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
