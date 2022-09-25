package dotodag

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"
)

// env:prod AND resource_name NOT IN (!resource_name:A, !resource_name:B)
// env:shop.ist AND availability-zone IN (us-east-1a, us-east-1b, us-east4-b)} by {availability-zone}
// avg:mymetric{env:prod AND resource_name NOT IN (!resource_name:A, !resource_name:B)}
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
			query:    "avg(last_5m):max:system.disk.in_use{*} by {host} > 1",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test simple monitor with float",
			query:    "avg(last_5m):max:system.disk.in_use{*} by {host} > 1.2",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test evaluation double-digit evaluation window",
			query:    "avg(last_15m):max:system.disk.in_use{*} by {host} > 1.2",
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
