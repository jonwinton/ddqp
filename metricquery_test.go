package dotodag

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"
)

func TestBool_Capture(t *testing.T) {
	type args struct {
		v []string
	}
	tests := []struct {
		name    string
		b       *Bool
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Capture(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Bool.Capture() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_MetricQuery(t *testing.T) {
	parser := participle.MustBuild[MetricQuery](participle.Unquote())

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
