package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type MetricFilter struct {
	Pos lexer.Position

	Parameters []*Param `( @@ ( ("," | "AND" | "and" | "OR" | "or" ) @@ )* | "*" )?`
}

func (mf *MetricFilter) String() string {
	if len(mf.Parameters) == 0 {
		return ""
	}

	params := []string{}
	for _, v := range mf.Parameters {
		params = append(params, v.String())
	}

	return strings.Join(params, ",")
}

type Param struct {
	GroupedFilter *GroupedFilter `"(" @@ ")"`
	SimpleFilter  *SimpleFilter  `| @@`
}

func (p *Param) String() string {
	if p.GroupedFilter != nil {
		return p.GroupedFilter.String()
	}

	return p.SimpleFilter.String()
}

type SimpleFilter struct {
	Negative        bool             `@"!"?`
	FilterKey       string           `@Ident`
	FilterSeparator *FilterSeparator `@@`
	FilterValue     *FilterValue     `@@`
}

func (sf *SimpleFilter) String() string {
	// TODO: Negative
	base := fmt.Sprintf("%s%s%s", sf.FilterKey, sf.FilterSeparator.String(), sf.FilterValue.String())

	if sf.Negative {
		return fmt.Sprintf("!%s", base)
	}

	return base
}

type GroupedFilter struct {
	Parameters []*Param `( @@ ( ("," | "AND" | "and" | "OR" | "or" ) @@ )* | "*" )?`
}

func (gf *GroupedFilter) String() string {
	params := []string{}
	for _, v := range gf.Parameters {
		params = append(params, v.String())
	}

	return strings.Join(params, ",")
}

type FilterSeparator struct {
	Colon bool `@(":" `
	In    bool `| ("IN" | "in") `
	NotIn bool `| ("NOT" "IN" | "not" "in") )`
}

func (fs *FilterSeparator) String() string {
	if fs.Colon {
		return ":"
	}

	if fs.In {
		return " IN "
	}

	return " NOT IN "
}

type FilterKey struct {
	Negative bool   `@"!"?`
	Key      string `@Ident`
}

func (fk *FilterKey) String() string {
	if fk.Negative {
		return fmt.Sprintf("!%s", fk.Key)
	}

	return fk.Key
}

type FilterValue struct {
	SimpleValue *Value   `	@@`
	ListValue   []*Value `| ( "(" ( @@ ( "," @@ | "OR" @@ | "or" @@ )* )? ")" )?`
}

func (fv *FilterValue) String() string {
	if len(fv.ListValue) > 0 {
		strs := []string{}
		for _, v := range fv.ListValue {
			strs = append(strs, v.String())
		}
		return strings.Join(strs, ",")
	}

	return fv.SimpleValue.String()
}

type Value struct {
	Boolean    *Bool    `  @("true"|"false")`
	Identifier *string  `| "!"? @Ident ( @"." @Ident )*`
	Str        *string  `| @(String)`
	Number     *float64 `| @(Float|Int)`
}

func (v *Value) String() string {
	if v.Boolean != nil {
		return v.Boolean.String()
	}

	if v.Number != nil {
		return fmt.Sprintf("%g", *v.Number)
	}

	if v.Identifier != nil {
		return *v.Identifier
	}

	return *v.Str
}
