package ddqp

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	// This is the primary lexer for all parsers
	// nolint:govet
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?i)rem[^\n]*`},
		{"String", `"(\\"|[^"])*"`},
		{"SpaceAggregatorCondition", `v: v[<>=]*([0-9]*[.])?[0-9]+`},
		{"ComparisonOperator", `:>[=]?|:<[=]?|:~`},
		{"Ident", `[a-zA-Z0-9_][\w\d-\*\./]*`},
		{"Float", `[+-]?([0-9]*[.])?[0-9]+`},
		{"Int", `\d+`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?\/]|]`},
		{"EOL", `[\n\r]+`},
		{"whitespace", `[ \t]+`},
	})
)
