package dotodag

import "github.com/alecthomas/participle/v2/lexer"

type MetricMonitor struct {
	Pos lexer.Position

	Aggregation      string       `@Ident`
	EvaluationWindow string       `"(" @Ident ")" ":"`
	MetricQuery      *MetricQuery `@@`
	Comparator       string       `@( ">" | ">" "=" | "<" | "<" "=" )`
	Threshold        float64      `@(Int|Float)`
}
