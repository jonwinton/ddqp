package ddqp

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

	Parameters []*Param `( @@ ( ("," | "AND" | "and" | "OR" | "or" ) @@ )* | "*" )?`
}

type Param struct {
	GroupedFilter *GroupedFilter `"(" @@ ")"`
	SimpleFilter  *SimpleFilter  `| @@`
}

type SimpleFilter struct {
	Negative        bool             `@"!"?`
	FilterKey       string           `@Ident`
	FilterSeparator *FilterSeparator `@@`
	FilterValue     *FilterValue     `@@`
}

type GroupedFilter struct {
	Parameters []*Param `( @@ ( ("," | "AND" | "and" | "OR" | "or" ) @@ )* | "*" )?`
}

type FilterSeparator struct {
	Colon bool `@(":" `
	In    bool `| ("IN" | "in") `
	NotIn bool `| ("NOT" "IN" | "not" "in") )`
}

type FilterKey struct {
	Negative bool   `@"!"?`
	Key      string `@Ident`
}

type FilterValue struct {
	SimpleValue *Value   `	@@`
	ListValue   []*Value `| ( "(" ( @@ ( "," @@ | "OR" @@ | "or" @@ )* )? ")" )?`
}

type Value struct {
	Boolean    *Bool    `  @("true"|"false")`
	Identifier *string  `| "!"? @Ident ( @"." @Ident )*`
	String     *string  `| @(String)`
	Number     *float64 `| @(Float|Int)`
}
