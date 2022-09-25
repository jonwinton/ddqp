package dotodag

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func NewMetricFilterParser() *participle.Parser[MetricFilter] {
	return participle.MustBuild[MetricFilter](
		participle.Lexer(lex),
		participle.Unquote("String"),
	)
}

type MetricFilter struct {
	Pos lexer.Position

	Parameters []*Param `( @@ ( ("," | "AND" | "OR" ) @@ )* | "*" )?`
}

type Param struct {
	Negative  bool   `@"!"?`
	FilterKey string `@Ident`
	// Group           []*Param         `| "(" @@ ")" )`
	FilterSeparator *FilterSeparator `@@`
	FilterValue     *FilterValue     `@@`
}

// typeg

type FilterSeparator struct {
	Colon bool `@(":" `
	In    bool `| "IN"`
	NotIn bool `| "NOT" "IN")`
}

type FilterKey struct {
	Negative bool   `@"!"?`
	Key      string `@Ident`
}

type FilterValue struct {
	SimpleValue *Value   `	@@`
	ListValue   []*Value `| ( "(" ( @@ ( "," @@ | "OR" @@ )* )? ")" )?`
}

type Value struct {
	Boolean    *Bool    `  @("true"|"false")`
	Identifier *string  `| "!"? @Ident ( @"." @Ident )*`
	String     *string  `| @(String)`
	Number     *float64 `| @(Float|Int)`
}
