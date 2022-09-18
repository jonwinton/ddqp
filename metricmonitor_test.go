package dotodag

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"
)

func Test_MetricMonitor(t *testing.T) {
	parser := NewMetricMonitorParser()

	tests := []struct {
		name     string
		query    string
		wantErr  bool
		printAST bool // For debugging, can opt in to print AST
	}{
		{
			name:     "test simple monitor",
			query:    "avg(last_5m):max:system.disk.in_use{fo:bf} by {host} > 1",
			wantErr:  false,
			printAST: true,
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
