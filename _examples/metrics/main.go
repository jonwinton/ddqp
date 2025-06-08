// Example of using DDQP to parse and work with DataDog metric queries
package main

import (
	"fmt"
	"github.com/jonwinton/ddqp"
)

// This example demonstrates how to:
// 1. Parse a DataDog metric query into a structured object
// 2. Access components of the query (aggregator, metric name, filters, grouping)
// 3. Convert the query back to its string representation
// 4. Create a new query programmatically
func main() {
	// Create a metric query parser
	parser := ddqp.NewMetricQueryParser()

	// Example 1: Parse an existing query
	fmt.Println("=== Example 1: Parse a metric query ===")
	query := `sum:kubernetes.containers.state.terminated{reason:oomkilled} by {kube_cluster_name,kube_deployment}`
	parsed, err := parser.Parse(query)
	if err != nil {
		panic(err)
	}

	// Access the structured data
	fmt.Printf("Original Query: %s\n", query)
	fmt.Printf("Aggregator: %s\n", parsed.Query[0].Aggregator)
	fmt.Printf("Metric Name: %s\n", parsed.Query[0].MetricName)
	fmt.Printf("Filter String: %s\n", parsed.Query[0].Filters.String())
	fmt.Printf("Filter Key: %s\n", parsed.Query[0].Filters.Left.SimpleFilter.FilterKey)
	fmt.Printf("Filter Value: %s\n", parsed.Query[0].Filters.Left.SimpleFilter.FilterValue.SimpleValue.Identifier)
	fmt.Printf("Grouping: %v\n", parsed.Query[0].Grouping)

	// Example 2: Working with complex filters
	fmt.Println("\n=== Example 2: Working with complex filters ===")
	queryWithComplexFilter := `sum:system.cpu.user{env:prod AND (service:api OR service:web)} by {host}`
	parsedComplex, err := parser.Parse(queryWithComplexFilter)
	if err != nil {
		panic(err)
	}

	// Note that the filter structure is more complex and can be traversed programmatically
	fmt.Printf("Complex Filter Query: %s\n", queryWithComplexFilter)
	fmt.Printf("Parsed Query String Representation: %s\n", parsedComplex.String())

	// Example 3: Parsing queries with special operators
	fmt.Println("\n=== Example 3: Queries with special operators ===")
	queryWithOperators := `sum:http.requests.count{status:>400 AND latency:<500}`
	parsedOperators, err := parser.Parse(queryWithOperators)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Operator Query: %s\n", queryWithOperators)
	fmt.Printf("Parsed Operator Query: %s\n", parsedOperators.String())

	// Example 4: Create metric query programmatically
	// This is a simplified example. In real usage, you would set more fields.
	fmt.Println("\n=== Example 4: Create a metric query programmatically ===")
	newParsed, _ := parser.Parse(`avg:system.cpu.idle{*} by {host}`)
	
	// Modify some fields
	newParsed.Query[0].Aggregator = "max"
	newParsed.Query[0].MetricName = "system.memory.used"
	
	// Create a simple filter programmatically
	newFilter := &ddqp.MetricFilter{}
	newFilter.Left = &ddqp.Param{}
	newFilter.Left.SimpleFilter = &ddqp.SimpleFilter{}
	newFilter.Left.SimpleFilter.FilterKey = "env"
	newFilter.Left.SimpleFilter.FilterSeparator = &ddqp.FilterSeparator{Colon: true}
	newFilter.Left.SimpleFilter.FilterValue = &ddqp.FilterValue{}
	newFilter.Left.SimpleFilter.FilterValue.SimpleValue = &ddqp.Value{}
	identifier := "production"
	newFilter.Left.SimpleFilter.FilterValue.SimpleValue.Identifier = &identifier
	
	// Replace the filter
	newParsed.Query[0].Filters = newFilter
	
	fmt.Printf("Programmatically Created Query: %s\n", newParsed.String())
}
