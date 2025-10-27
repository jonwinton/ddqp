// Example of using DDQP to work with DataDog monitor queries
package main

import (
	"fmt"

	"github.com/jonwinton/ddqp"
)

// This example demonstrates working with DataDog monitor queries:
// 1. Parsing monitor queries with thresholds
// 2. Accessing monitor components (evaluation window, comparator, threshold)
// 3. Different types of monitors (min, max, avg)
// 4. Using complex filters in monitors
func main() {
	// Create a monitor parser
	parser := ddqp.NewMetricMonitorParser()

	// Example 1: Basic monitor query
	fmt.Println("=== Example 1: Basic Monitor Query ===")
	monitorQuery := `avg(last_5m):avg:system.cpu.user{env:prod} > 80`

	monitor, err := parser.Parse(monitorQuery)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Monitor Query: %s\n", monitorQuery)
	fmt.Printf("Aggregation: %s\n", monitor.Aggregation)
	fmt.Printf("Evaluation Window: %s\n", monitor.EvaluationWindow)
	fmt.Printf("Metric Name: %s\n", monitor.MetricQuery.Query[0].MetricName)
	fmt.Printf("Filters: %s\n", monitor.MetricQuery.Query[0].Filters.String())
	fmt.Printf("Comparator: %s\n", monitor.Comparator)
	fmt.Printf("Threshold: %g\n", monitor.Threshold)

	// Example 2: Different monitor types
	fmt.Println("\n=== Example 2: Different Monitor Types ===")

	// Min monitor
	minMonitor := `min(last_10m):min:system.cpu.idle{env:production} < 10`
	minParsed, _ := parser.Parse(minMonitor)

	// Max monitor
	maxMonitor := `max(last_30m):max:system.memory.used{service:database} > 90`
	maxParsed, _ := parser.Parse(maxMonitor)

	// Sum monitor
	sumMonitor := `sum(last_1h):sum:system.io.await{service:api} >= 100`
	sumParsed, _ := parser.Parse(sumMonitor)

	fmt.Printf("Min Monitor: %s\n    -> %s\n", minMonitor, minParsed.String())
	fmt.Printf("Max Monitor: %s\n    -> %s\n", maxMonitor, maxParsed.String())
	fmt.Printf("Sum Monitor: %s\n    -> %s\n", sumMonitor, sumParsed.String())

	// Example 3: Monitor with complex filter
	fmt.Println("\n=== Example 3: Monitor with Complex Filter ===")
	complexMonitor := `avg(last_15m):avg:network.tcp.retransmit{env:prod AND (region:us-east OR region:us-west)} > 50`

	complexParsed, err := parser.Parse(complexMonitor)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Complex Monitor: %s\n", complexMonitor)
	fmt.Printf("Parsed Monitor: %s\n", complexParsed.String())

	// Example 4: Monitor with grouping
	fmt.Println("\n=== Example 4: Monitor with Grouping ===")
	groupingMonitor := `min(last_5m):min:system.load.1{env:production} by {host,availability-zone,cluster} > 4`

	groupingParsed, err := parser.Parse(groupingMonitor)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grouping Monitor: %s\n", groupingMonitor)
	fmt.Printf("Metric Name: %s\n", groupingParsed.MetricQuery.Query[0].MetricName)
	fmt.Printf("Grouping: %v\n", groupingParsed.MetricQuery.Query[0].Grouping)

	// Example 5: Creating a monitor programmatically
	fmt.Println("\n=== Example 5: Creating a Monitor Programmatically ===")

	// Create a new metric query structure
	metricQuery := &ddqp.MetricQuery{}
	metricQuery.Query = []*ddqp.Query{{
		Aggregator: "avg",
		MetricName: "system.disk.in_use",
		Filters: &ddqp.MetricFilter{
			Left: &ddqp.Param{
				SimpleFilter: &ddqp.SimpleFilter{
					FilterKey:       "service",
					FilterSeparator: &ddqp.FilterSeparator{Colon: true},
					FilterValue: &ddqp.FilterValue{
						SimpleValue: &ddqp.Value{
							Identifier: func() *string { s := "database"; return &s }(),
						},
					},
				},
			},
		},
		By:       "by",
		Grouping: []string{"host"},
	}}

	// Create the monitor
	newMonitor := &ddqp.MetricMonitor{
		Aggregation:      "max",
		EvaluationWindow: "last_15m",
		MetricQuery:      metricQuery,
		Comparator:       ">",
		Threshold:        95,
	}

	fmt.Printf("Programmatically Created Monitor: %s\n", newMonitor.String())
}
