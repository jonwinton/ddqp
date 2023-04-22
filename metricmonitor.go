package ddqp

import (
	"fmt"
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
	Threshold        float64      `@(Ident)`
}

// String returns the string representation of the metric monitor.
func (mm *MetricMonitor) String() string {
	return fmt.Sprintf("%s(%s):%s %s %g", mm.Aggregation, mm.EvaluationWindow, mm.MetricQuery.String(), mm.Comparator, mm.Threshold)
}

// NewMetricMonitorParser returns a Parser which is capable of interpretting
// a metric query.
func NewMetricMonitorParser() *MetricMonitorParser {
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
func (mmp *MetricMonitorParser) Parse(query string) (*MetricMonitor, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(query, "\n", "")
	// return the raw parsed outpu
	return mmp.parser.ParseString("", sanitized)
}
