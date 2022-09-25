package dotodag

import "github.com/alecthomas/participle/v2/lexer"

// https://docs.datadoghq.com/metrics/#anatomy-of-a-metric-query
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
