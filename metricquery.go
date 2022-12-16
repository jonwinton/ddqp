package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type MetricQuery struct {
	Pos lexer.Position

	Query []*Query `parser:"@@"`
}

func (mq *MetricQuery) String() string {
	return mq.Query[0].String()
}

type Query struct {
	Pos lexer.Position

	Aggregator string        `parser:"@Ident ':'"`
	MetricName string        `parser:"@Ident( @'.' @Ident)*"`
	Filters    *MetricFilter `"{" @@ "}"`
	By         string        `parser:"Ident?"`
	Grouping   []string      `parser:"'{'? ( @Ident ( ',' @Ident )* )? '}'?"`
	Function   []*Function   `parser:"( @@ ( '.' @@ )* )?"`
}

func (q *Query) String() string {
	base := fmt.Sprintf("%s:%s{%s}", q.Aggregator, q.MetricName, q.Filters.String())

	if len(q.Grouping) > 0 {
		base = fmt.Sprintf("%s by {%s}", base, strings.Join(q.Grouping, ","))
	}

	if len(q.Function) > 0 {
		funcs := []string{}
		for _, v := range q.Function {
			funcs = append(funcs, v.String())
		}
		return fmt.Sprintf("%s.%s", base, strings.Join(funcs, "."))
	}

	return base
}

type Function struct {
	Name string   `"." @Ident`
	Args []*Value `"(" ( @@ ( "," @@ )* )? ")"`
}

func (f *Function) String() string {
	args := []string{}
	for _, v := range f.Args {
		args = append(args, v.String())
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ","))
}

type Bool bool

func (b *Bool) Capture(v []string) error { *b = v[0] == "true"; return nil }
func (b *Bool) String() string           { return fmt.Sprintf("%v", *b) }

// NewMetricQueryParser returns a Parser which is capable of interpretting
// a metric query.
func NewMetricQueryParser() *MetricQueryParser {
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
