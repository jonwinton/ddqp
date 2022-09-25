package dotodag

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"
)

func Test_MetricQuery(t *testing.T) {
	parser := NewMetricQueryParser()

	tests := []struct {
		name     string
		query    string
		wantErr  bool
		printAST bool // For debugging, can opt in to print AST
	}{
		// Simple passing example. Guaranteed to have all parts of a query
		{
			name:     "simple query",
			query:    "sum:namespace.metric.name{foo:bar,baz:bang} by {foo,bar}",
			wantErr:  false,
			printAST: false,
		},
		// Simple failing example. Guaranteed to fail because missing the aggregator
		{
			name:     "fail due to no aggregator",
			query:    "namespace.metric.name{foo:bar,baz:bang} by {foo,bar}",
			wantErr:  true,
			printAST: false,
		},
		{
			name:     "filter by asterisk",
			query:    "sum:namespace.metric.name{*} by {foo,bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "filer by partial asterisk",
			query:    "sum:namespace.metric.name{foo:bar-*} by {foo,bar}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test underscores in metric name",
			query:    "sum:namespace.metric_name{foo:bar} by {baz}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test hyphens in filters and groupings",
			query:    "sum:prometheus_metric_source{foo:bar-bar,baz:bang} by {fizz-buzz,bang}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numbers in the metric name",
			query:    "sum:prometheus_metric_source_1{foo:bar-bar,baz:bang} by {fizz-buzz,bang}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numbers in the filters and grouping",
			query:    "sum:prometheus_metric_source_1{foo:bar-bar-1,baz:bang_2} by {fizz-buzz3,bang}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test docs example query",
			query:    "avg:system.cpu.user{env:staging AND (availability-zone:us-east-1a OR availability-zone:us-east-1c)} by {availability-zone}",
			wantErr:  false,
			printAST: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.ParseString("", tt.query)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err)
			}

			if tt.printAST {
				repr.Println(ast)
			}
		})
	}
}
