package ddqp

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// This is the primary lexer for all parsers
// nolint:govet
var lex = lexer.MustSimple([]lexer.SimpleRule{
	{"Comment", `(?i)rem[^\n]*`},
	{"String", `"(\\"|[^"])*"`},
	{"SpaceAggregatorCondition", `v: v[<>=]*([0-9]*[.])?[0-9]+`},
	{"Operator", `[<>=*\/+-]+`},
	{"OperatorValue", `[<>=]+([0-9]*[.])?[0-9]+`},
	{"RangeValue", `\[[0-9]+ TO [0-9]+\]`},
	{"Ident", `[a-zA-Z0-9_\*][\w\d-\*\./\?]*`},
	{"Float", `[+-]?([0-9]*[.])?[0-9]+`},
	{"Int", `\d+`},
	{"Attribute", `@[a-zA-Z0-9_\*][\w\d-\*\./\?]*`},
	{"BooleanOp", `(?i)(AND|OR|NOT)`},
	{"CIDR", `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}`},
	{"RangeOp", `(?i)TO`},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?\/]|]`},
	{"EOL", `[\n\r]+`},
	{"whitespace", `[ \t]+`},
})
