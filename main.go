package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// https://docs.datadoghq.com/metrics/#anatomy-of-a-metric-query
type MetricQuery struct {
	Pos lexer.Position

	Query []*Query `@@`
}

type Query struct {
	Pos lexer.Position

	Aggregator string      `@Ident ":"`
	MetricName string      `@Ident( @"." @Ident)*`
	Filters    []*Filter   `"{" ( @@ ( "," @@ )* )? "}"`
	Function   []*Function `( @@ ( "." @@ )* )?`
	By         string      `Ident`
	Grouping   []string    `"{" ( @Ident ( "," @Ident )* )? "}"`
}
type Filter struct {
	Key   string `@Ident ":"`
	Value string `@(String|Ident)`
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
	String     *string  `| @(String|Char|RawString)`
	Number     *float64 `| @(Float|Int)`
}

func main() {
	parser := participle.MustBuild[MetricQuery](
		participle.Unquote("String"),
	)

	query, err := parser.ParseString("", `sum:kubernetes.containers.state.terminated{reason:oomkilled} by    {kube_cluster_name,kube_deployment}`)
	if err != nil {
		panic(err)
	}
	repr.Println(query)
}
