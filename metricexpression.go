package ddqp

import (
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
	Number        *float64     `  @(Float|Int)`
	MetricQuery   *MetricQuery `| @@`
	Subexpression *Expression  `| "(" @@ ")"`
}

type Factor struct {
	Base *ExprValue `@@`
	// Exponent *ExprValue `( "^" @@ )?`
}

type OpFactor struct {
	Operator Operator `@("*" | "/")`
	Factor   *Factor  `@@`
}

type Term struct {
	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Operator Operator `@("+" | "-")`
	Term     *Term    `@@`
}

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type MetricExpression struct {
	Pos lexer.Position

	Left  *Term     `@@`
	Right []*OpTerm `@@*`
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
