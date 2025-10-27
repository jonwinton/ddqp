package ddqp

import (
	"strings"
)

type GenericParser struct{}

func NewGenericParser() *GenericParser {
	return &GenericParser{}
}

type GenericQuery struct {
	MetricExpression *MetricExpression
	MetricQuery      *MetricQuery
}

// Parse sanitizes the query string and returns the AST and any error.
func (gp *GenericParser) Parse(query string) (*GenericQuery, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(query, "\n", "")

	// Prefer MetricQuery parsing first because '*' '-' '/' are valid inside identifiers/filters
	mqp := NewMetricQueryParser()
	if metricQuery, err := mqp.Parse(sanitized); err == nil {
		return &GenericQuery{MetricQuery: metricQuery}, nil
	}

	// Fallback to MetricExpression
	mep := NewMetricExpressionParser()
	metricExpression, err := mep.Parse(sanitized)
	if err != nil {
		return nil, err
	}
	return &GenericQuery{MetricExpression: metricExpression}, nil
}

func (gq *GenericQuery) String() string {
	if gq.MetricExpression != nil {
		return gq.MetricExpression.String()
	} else if gq.MetricQuery != nil {
		return gq.MetricQuery.String()
	}
	return ""
}
