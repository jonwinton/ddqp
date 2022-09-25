package ddqp

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type MetricQuery struct {
	Pos lexer.Position

	Query []*Query `@@`
}

type Query struct {
	Pos lexer.Position

	Aggregator string        `@Ident ":"`
	MetricName string        `@Ident( @"." @Ident)*`
	Filters    *MetricFilter `"{" @@ "}"`
	Function   []*Function   `( @@ ( "." @@ )* )?`
	By         string        `Ident`
	Grouping   []string      `"{" ( @Ident ( "," @Ident )* )? "}"`
}
type Filter struct {
	Key   string `@Ident ":"`
	Value string `@Ident`
}
type Function struct {
	Name string          `"." @Ident`
	Args []*FunctionArgs `"(" ( @@ ( "," @@ )* )? ")"`
}

type Bool bool

func (b *Bool) Capture(v []string) error { *b = v[0] == "true"; return nil }

type FunctionArgs struct {
	Boolean    *Bool    `  @("true"|"false")`
	Identifier *string  `| @Ident ( @"." @Ident )*`
	String     *string  `| @(String)`
	Number     *float64 `| @(Float|Int)`
}

// NewMetricQueryParser returns a Parser which is capable of interpretting
// a metric query.
func NewMetricQueryParser() Parser {
	mqp := &MetricQueryParser{
		parser: participle.MustBuild[MetricQuery](
			participle.Lexer(lex),
			participle.Unquote("String"),
		),
	}

	return mqp
}

// MetricQueryParser is parser returned when calling NewMetricQueryParser.
type MetricQueryParser struct {
	parser *participle.Parser[MetricQuery]
}

// Parse sanitizes the query string and returns the AST and any error.
func (mqp *MetricQueryParser) Parse(query string) (*MetricQuery, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(query, "\n", "")
	// return the raw parsed outpu
	return mqp.parser.ParseString("", sanitized)
}
