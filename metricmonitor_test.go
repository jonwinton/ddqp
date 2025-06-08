package ddqp

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
		{
			name:     "test less than operator",
			query:    "min(last_10m):min:system.cpu.idle{env:production} by {host} < 10",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test greater than or equal operator",
			query:    "avg(last_30m):sum:errors.count{service:api} by {endpoint} > 500",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test less than or equal operator",
			query:    "max(last_1h):avg:system.memory.free{role:database} < 100",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test complex filter in monitor",
			query:    "avg(last_15m):max:network.tcp.retransmit{env:prod AND (region:us-east OR region:us-west)} > 50",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test complex grouping in monitor",
			query:    "min(last_5m):avg:system.load.1{env:production} by {host,availability-zone,cluster} > 4",
			wantErr:  false,
			printAST: false,
		},
		{
			name:     "test numeric threshold with decimal places",
			query:    "avg(last_5m):sum:system.io.await{service:database} > 25.75",
			wantErr:  false,
			printAST: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.Parse(tt.query)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err)
			}

			if tt.printAST {
				repr.Println(ast)
			}

			// test re-stringification
			restring := ast.String()
			require.Equal(t, tt.query, restring)
		})
	}
}
