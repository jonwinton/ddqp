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
			query:    "sum:kubernetes.containers.state.terminated{reason:oomkilled} by {kube_cluster_name,kube_deployment}",
			wantErr:  false,
			printAST: false,
		},
		// Simple failing example. Guaranteed to fail because missing the aggregator
		{
			name:     "fail due to no aggregator",
			query:    "kubernetes.containers.state.terminated{reason:oomkilled} by {kube_cluster_name,kube_deployment}",
			wantErr:  true,
			printAST: false,
		},
		{
			name:     "test underscores in metric name",
			query:    "sum:prometheus_metric_source{foo:bar} by {baz}",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test hyphens in filters name",
			query:    "sum:prometheus_metric_source{foo:bar-bar} by {baz}",
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
