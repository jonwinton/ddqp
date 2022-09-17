package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// https://docs.datadoghq.com/metrics/#anatomy-of-a-metric-query
type MetricQuery struct {
	Pos lexer.Position

	SpaceAggregator string `@Ident ":"`
	MetricName      string `@Ident {`
}

func main() {
	tomlLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"DateTime", `\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?(-\d\d:\d\d)?`},
		{"Date", `\d\d\d\d-\d\d-\d\d`},
		{"Time", `\d\d:\d\d:\d\d(\.\d+)?`},
		{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*`},
		{"String", `"[^"]*"`},
		{"Number", `[-+]?[.0-9]+\b`},
		{"Punct", `\[|]|[-!()+/*=,]`},
		{"comment", `#[^\n]+`},
		{"whitespace", `\s+`},
	})
	parser := participle.MustBuild[MetricQuery](
		participle.Lexer(tomlLexer),
		participle.Unquote("String"),
	)

	fmt.Println(parser)
	query, err := parser.ParseString("", "avg:system.disk.free{*}.rollup(avg, 60)")
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	log.Println("hi")
	fmt.Println("hi")
}
