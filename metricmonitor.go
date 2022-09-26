package ddqp

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type MetricMonitor struct {
	Pos lexer.Position

	Aggregation      string       `@Ident`
	EvaluationWindow string       `"(" @Ident ")" ":"`
	MetricQuery      *MetricQuery `@@`
	Comparator       string       `@( ">" | ">" "=" | "<" | "<" "=" )`
	Threshold        float64      `@(Int|Float)`
}

// NewMetricMonitorParser returns a Parser which is capable of interpretting
// a metric query.
func NewMetricMonitorParser() Parser {
	mmp := &MetricMonitorParser{
		parser: participle.MustBuild[MetricMonitor](
			participle.Lexer(lex),
			participle.Unquote("String"),
		),
	}

	return mmp
}

// MetricMonitorParser is parser returned when calling NewMetricMonitorParser.
type MetricMonitorParser struct {
	parser *participle.Parser[MetricMonitor]
}

// Parse sanitizes the query string and returns the AST and any error.
func (mmp *MetricMonitorParser) Parse(query string) (ParsedResponse, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(query, "\n", "")
	// return the raw parsed outpu
	return mmp.parser.ParseString("", sanitized)
}
