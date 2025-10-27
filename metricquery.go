package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type MetricQuery struct {
	Pos               lexer.Position
	Query             *Query             `parser:"@@"`
	AggregatorFuction *AggregatorFuction `parser:"| @@"`
}

type AggregatorFuction struct {
	Pos  lexer.Position
	Name string       `parser:"@Ident '('"`
	Body *MetricQuery `parser:"@@"` // Query or AggregatorFuction
	Args []*Value     `parser:"( ',' @@ )* ')'"`
}

func (mq *MetricQuery) String() string {
	if mq.Query != nil {
		return mq.Query.String()
	}
	return mq.AggregatorFuction.String()
}

// String prints wrapper(body, arg1, arg2, ...), allowing nested wrappers as body.
func (w *AggregatorFuction) String() string {
	args := []string{}
	for _, v := range w.Args {
		args = append(args, v.String())
	}
	argTail := ""
	if len(args) > 0 {
		argTail = ", " + strings.Join(args, ", ")
	}
	return fmt.Sprintf("%s(%s%s)", w.Name, w.Body.String(), argTail)
}

type Query struct {
	Pos lexer.Position

	Aggregator *Aggregator   `parser:"@@?"`
	MetricName string        `parser:"@Ident( @'.' @Ident)*"`
	Filters    *MetricFilter `parser:"'{' @@ '}'"`
	By         string        `parser:"('by')?"`
	Grouping   []string      `parser:"( '{' ( @(Ident|'*') ( ',' @(Ident|'*') )* ) '}' )?"`
	Function   []*Function   `parser:"( '.' @@ ( '.' @@ )* )?"`
}

type Aggregator struct {
	Pos                       lexer.Position
	Name                      string `parser:"@Ident"`
	SpaceAggregationCondition string `parser:"( '(' @SpaceAggregatorCondition ')' )?"`
	Separator                 string `parser:"':'"`
}

func (q *Query) String() string {
	base := ""
	if q.Aggregator != nil {
		base = q.Aggregator.Name
		if q.Aggregator.SpaceAggregationCondition != "" {
			base = fmt.Sprintf("%s(%s)", base, q.Aggregator.SpaceAggregationCondition)
		}
		base = fmt.Sprintf("%s:", base)
	}

	base = fmt.Sprintf("%s%s{%s}", base, q.MetricName, q.Filters.String())

	if len(q.Grouping) > 0 {
		base = fmt.Sprintf("%s by {%s}", base, strings.Join(q.Grouping, ","))
	}

	if len(q.Function) > 0 {
		funcs := []string{}
		for _, v := range q.Function {
			funcs = append(funcs, v.String())
		}
		base = fmt.Sprintf("%s.%s", base, strings.Join(funcs, "."))
	}

	return base
}

type Function struct {
	Name string   `parser:"@Ident"`
	Args []*Value `parser:"'(' ( @@ ( ',' @@ )* )? ')'"`
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
