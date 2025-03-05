package ddqp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryFilterParser(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
		wantErr  bool
	}{
		// Basic filters
		{
			name:     "simple tag filter",
			query:    "env:prod",
			expected: "env:prod",
			wantErr:  false,
		},
		{
			name:     "full-text search",
			query:    "*:hello",
			expected: "*:hello",
			wantErr:  false,
		},
		{
			name:     "full-text search with exact match",
			query:    "*:\"hello world\"",
			expected: "*:\"hello world\"",
			wantErr:  false,
		},

		// Boolean expressions
		{
			name:     "boolean expression with AND",
			query:    "env:prod AND service:web",
			expected: "env:prod AND service:web",
			wantErr:  false,
		},
		{
			name:     "boolean expression with OR",
			query:    "env:prod OR env:staging",
			expected: "env:prod OR env:staging",
			wantErr:  false,
		},

		// Negation
		{
			name:     "negated tag",
			query:    "-version:beta",
			expected: "-version:beta",
			wantErr:  false,
		},

		// Grouping
		{
			name:     "grouped expression",
			query:    "(env:prod OR env:staging)",
			expected: "(env:prod OR env:staging)",
			wantErr:  false,
		},
		{
			name:     "complex boolean expression",
			query:    "(env:prod OR env:staging) AND -version:beta",
			expected: "(env:prod OR env:staging) AND -version:beta",
			wantErr:  false,
		},

		// Wildcards
		{
			name:     "wildcard in tag value",
			query:    "service:web*",
			expected: "service:web*",
			wantErr:  false,
		},
		{
			name:     "wildcard at beginning of tag value",
			query:    "service:*web",
			expected: "service:*web",
			wantErr:  false,
		},

		// Attribute filters
		{
			name:     "simple attribute filter",
			query:    "@http.url:www.datadoghq.com",
			expected: "@http.url:www.datadoghq.com",
			wantErr:  false,
		},
		{
			name:     "attribute filter with operator",
			query:    "@http.response_time:>100",
			expected: "@http.response_time:>100",
			wantErr:  true,
		},
		{
			name:     "attribute filter with range",
			query:    "@http.status_code:[400 TO 499]",
			expected: "@http.status_code:[400 TO 499]",
			wantErr:  false,
		},

		// Special cases
		{
			name:     "search with CIDR notation",
			query:    "@ip_address:192.168.1.0/24",
			expected: "@ip_address:192.168.1.0/24",
			wantErr:  false,
		},
		{
			name:     "search with question mark wildcard",
			query:    "@my_attribute:hello?world",
			expected: "@my_attribute:hello?world",
			wantErr:  false,
		},
	}

	parser := NewQueryFilterParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.query)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}
