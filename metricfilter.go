package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type MetricFilter struct {
	Pos lexer.Position

	Left       *Param   `(@@ | "*" )`
	Parameters []*Param `( @@* )`
}

func (mf *MetricFilter) String() string {
	params := []string{
		mf.Left.String(),
	}
	for _, v := range mf.Parameters {
		params = append(params, v.String())
	}

	return strings.Join(params, "")
}

type Param struct {
	GroupedFilter *GroupedFilter        ` "(" @@ ")"`
	Separator     *FilterValueSeparator `| @@`
	SimpleFilter  *SimpleFilter         `| @@`
	Asterisk      bool                  `| @"*"`
}

func (p *Param) String() string {
	if p.Separator != nil {
		return p.Separator.String()
	}

	if p.GroupedFilter != nil {
		return p.GroupedFilter.String()
	}

	if p.Asterisk {
		return "*"
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
	Parameters []*Param `( @@* | "*" )?`
}

func (gf *GroupedFilter) String() string {
	params := []string{}
	for _, v := range gf.Parameters {
		params = append(params, v.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(params, ""))
}

type FilterSeparator struct {
	Colon        bool `@":"`
	GreaterThan  bool `| @":>"`
	LessThan     bool `| @":<"`
	GreaterEqual bool `| @":>="`
	LessEqual    bool `| @":<="`
	Regex        bool `| @":~"`
	In           bool `| @("IN" | "in") `
	NotIn        bool `| @("NOT" "IN" | "not" "in")`
	AndNot       bool `| @("AND" "NOT" | "and" "not")`
}

func (fs *FilterSeparator) String() string {
	if fs.Colon {
		return ":"
	}

	if fs.GreaterThan {
		return ":>"
	}

	if fs.LessThan {
		return ":<"
	}

	if fs.GreaterEqual {
		return ":>="
	}

	if fs.LessEqual {
		return ":<="
	}

	if fs.Regex {
		return ":~"
	}

	if fs.In {
		return " IN "
	}

	if fs.AndNot {
		return " AND NOT "
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
	SimpleValue *Value   `@@`
	ListValue   []*Value `| ( "(" @@* ")" )?`
}

func (fv *FilterValue) String() string {
	if len(fv.ListValue) > 0 {
		strs := []string{}
		for _, v := range fv.ListValue {
			strs = append(strs, v.String())
		}
		return fmt.Sprintf("(%s)", strings.Join(strs, ""))
	}

	return fv.SimpleValue.String()
}

type Value struct {
	Separator  *FilterValueSeparator ` @@`
	Boolean    *Bool                 `|  @("true"|"false")`
	Identifier *string               `| "!"? @Ident ( @"." @Ident )*`
	Str        *string               `| @(String)`
	Number     *float64              `| @(Float|Int)`
	Wildcard   *string               `| @(FilterIdent|'*')`
}

func (v *Value) String() string {
	if v.Boolean != nil {
		return v.Boolean.String()
	}

	if v.Number != nil {
		return formatFloatNoExp(*v.Number)
	}

	if v.Identifier != nil {
		return *v.Identifier
	}

	if v.Separator != nil {
		return v.Separator.String()
	}

	if v.Wildcard != nil {
		return *v.Wildcard
	}

	return *v.Str
}

type FilterValueSeparator struct {
	Comma  bool ` @","`
	AndNot bool `| @("AND" "NOT" | "and" "not")`
	And    bool `| @("AND" | "and")`
	Or     bool `| @("OR" | "or")`
	In     bool `| @("IN" | "in")`
}

func (fvs *FilterValueSeparator) String() string {
	if fvs.Comma {
		return ", "
	}

	if fvs.And {
		return " AND "
	}

	if fvs.Or {
		return " OR "
	}

	if fvs.AndNot {
		return " AND NOT "
	}

	return " IN "
}
