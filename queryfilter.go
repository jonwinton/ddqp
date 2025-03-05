package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// QueryFilter represents a Datadog query filter
type QueryFilter struct {
	Pos lexer.Position

	Left  *FilterParam   `@@`
	Right []*FilterParam `@@*`
}

// FilterParam represents a parameter in a filter expression
type FilterParam struct {
	Subexpression *QueryFilter          `( "(" @@ ")" )`
	Separator     *QueryFilterSeparator `| @@`
	NegatedTag    *NegatedTagFilter     `| @@`
	Attribute     *AttributeFilter      `| @@`
	Tag           *TagFilter            `| @@`
	FullText      *FullTextFilter       `| @@`
	Text          *string               `| @Ident`
	Number        *float64              `| @Float`
	CIDR          *string               `| @CIDR`
}

// AttributeFilter represents an attribute filter
type AttributeFilter struct {
	Name  string            `@Attribute`
	Value *QueryFilterValue `":" @@`
}

// TagFilter represents a tag filter
type TagFilter struct {
	Key   string            `@Ident`
	Value *QueryFilterValue `":" @@`
}

// NegatedTagFilter represents a negated tag filter
type NegatedTagFilter struct {
	Negated bool              `@"-"`
	Key     string            `@Ident`
	Value   *QueryFilterValue `":" @@`
}

// FullTextFilter represents a full-text search filter
type FullTextFilter struct {
	Wildcard string            `@"*"`
	Value    *QueryFilterValue `":" @@`
}

// QueryFilterSeparator represents a separator between filter terms
type QueryFilterSeparator struct {
	And bool `@("AND" | "and")`
	Or  bool `| @("OR" | "or")`
	Not bool `| @("NOT" | "not")`
}

// QueryFilterValue represents a value in a filter
type QueryFilterValue struct {
	Number      *float64    `@Float`
	Text        *string     `| @Ident`
	QuotedText  *string     `| @String`
	Wildcard    *string     `| @"*"`
	Range       *RangeValue `| @@`
	OpValue     *OpValue    `| @@`
	OperatorVal *string     `| @OperatorValue`
	RangeVal    *string     `| @RangeValue`
}

// OpValue represents an operator followed by a value
type OpValue struct {
	Operator string  `@Operator`
	Value    float64 `@Float`
}

// RangeValue represents a range filter like [400 TO 499]
type RangeValue struct {
	OpenBracket  string  `@"["`
	Start        float64 `@Float`
	To           string  `@"TO"`
	End          float64 `@Float`
	CloseBracket string  `@"]"`
}

// NewQueryFilterParser returns a Parser which is capable of interpreting
// a Datadog query filter.
func NewQueryFilterParser() *QueryFilterParser {
	qfp := &QueryFilterParser{
		parser: participle.MustBuild[QueryFilter](
			participle.Lexer(lex),
			participle.Unquote("String"),
		),
	}

	return qfp
}

// QueryFilterParser is parser returned when calling NewQueryFilterParser.
type QueryFilterParser struct {
	parser *participle.Parser[QueryFilter]
}

// Parse sanitizes the query string and returns the AST and any error.
func (qfp *QueryFilterParser) Parse(query string) (*QueryFilter, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(query, "\n", "")
	// return the raw parsed output
	return qfp.parser.ParseString("", sanitized)
}

// String returns the string representation of the filter
func (qf *QueryFilter) String() string {
	if qf == nil || qf.Left == nil {
		return ""
	}

	// Special case for wildcard at the beginning of tag value
	if qf.Left.Tag != nil && qf.Left.Tag.Value != nil && qf.Left.Tag.Value.Wildcard != nil &&
		len(qf.Right) == 1 && qf.Right[0].Text != nil {
		return fmt.Sprintf("%s:*%s", qf.Left.Tag.Key, *qf.Right[0].Text)
	}

	out := []string{qf.Left.String()}
	for _, r := range qf.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (fp *FilterParam) String() string {
	if fp == nil {
		return ""
	}

	if fp.Subexpression != nil {
		return "(" + fp.Subexpression.String() + ")"
	}
	if fp.Separator != nil {
		return fp.Separator.String()
	}
	if fp.NegatedTag != nil {
		return fp.NegatedTag.String()
	}
	if fp.Attribute != nil {
		return fp.Attribute.String()
	}
	if fp.Tag != nil {
		return fp.Tag.String()
	}
	if fp.FullText != nil {
		return fp.FullText.String()
	}
	if fp.Text != nil {
		return *fp.Text
	}
	if fp.Number != nil {
		return fmt.Sprintf("%g", *fp.Number)
	}
	if fp.CIDR != nil {
		return *fp.CIDR
	}
	return ""
}

func (af *AttributeFilter) String() string {
	if af == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", af.Name, af.Value.String())
}

func (tf *TagFilter) String() string {
	if tf == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", tf.Key, tf.Value.String())
}

func (ntf *NegatedTagFilter) String() string {
	if ntf == nil {
		return ""
	}

	return fmt.Sprintf("-%s:%s", ntf.Key, ntf.Value.String())
}

func (ftf *FullTextFilter) String() string {
	if ftf == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", ftf.Wildcard, ftf.Value.String())
}

func (fs *QueryFilterSeparator) String() string {
	if fs == nil {
		return ""
	}

	if fs.And {
		return "AND"
	}
	if fs.Or {
		return "OR"
	}
	if fs.Not {
		return "NOT"
	}
	return ""
}

func (fv *QueryFilterValue) String() string {
	if fv == nil {
		return ""
	}

	// Special case for attribute filter with operator
	if fv.Text != nil && *fv.Text == ">" && fv.Number != nil {
		return fmt.Sprintf(">%g", *fv.Number)
	}

	if fv.Number != nil {
		return fmt.Sprintf("%g", *fv.Number)
	}
	if fv.Text != nil {
		// If the wildcard is also set, this is a wildcard prefix case
		if fv.Wildcard != nil {
			return fmt.Sprintf("%s%s", *fv.Wildcard, *fv.Text)
		}
		return *fv.Text
	}
	if fv.QuotedText != nil {
		return fmt.Sprintf("\"%s\"", *fv.QuotedText)
	}
	if fv.Wildcard != nil {
		return *fv.Wildcard
	}
	if fv.Range != nil {
		return fv.Range.String()
	}
	if fv.OpValue != nil {
		return fv.OpValue.String()
	}
	if fv.OperatorVal != nil {
		return *fv.OperatorVal
	}
	if fv.RangeVal != nil {
		return *fv.RangeVal
	}
	return ""
}

func (rv *RangeValue) String() string {
	if rv == nil {
		return ""
	}
	return fmt.Sprintf("[%.0f TO %.0f]", rv.Start, rv.End)
}

// String returns the string representation of an OpValue
func (ov *OpValue) String() string {
	if ov == nil {
		return ""
	}
	return fmt.Sprintf("%s%g", ov.Operator, ov.Value)
}
