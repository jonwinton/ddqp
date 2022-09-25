package ddqp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	// nolint:govet
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?i)rem[^\n]*`},
		{"String", `"(\\"|[^"])*"`},
		{"Float", `[+-]?([0-9]*[.])?[0-9]+`},
		{"Int", `\d+`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?\/]|]`},
		{"Ident", `[a-zA-Z0-9_][\w\d-\*]*`},
		{"EOL", `[\n\r]+`},
		{"whitespace", `[ \t]+`},
	})
)

func NewMetricQueryParser() *participle.Parser[MetricQuery] {
	return participle.MustBuild[MetricQuery](
		participle.Lexer(lex),
		participle.Unquote("String"),
	)
}

func NewMetricMonitorParser() *participle.Parser[MetricMonitor] {
	return participle.MustBuild[MetricMonitor](
		participle.Lexer(lex),
		participle.Unquote("String"),
	)
}
