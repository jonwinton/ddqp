// Example of using DDQP to work with DataDog filters
package main

import (
	"fmt"
	"github.com/jonwinton/ddqp"
)

// This example demonstrates advanced techniques for working with filters in DataDog queries:
// 1. Creating and parsing complex boolean expressions (AND, OR, NOT)
// 2. Working with comparison operators (>, <, >=, <=)
// 3. Using regex filters
// 4. Building filters programmatically
func main() {
	parser := ddqp.NewMetricQueryParser()

	// Example 1: Complex boolean logic in filters
	fmt.Println("=== Example 1: Complex Boolean Logic in Filters ===")
	query := `sum:http.requests{env:prod AND (service:api OR service:web) AND NOT status:500}`
	
	parsed, err := parser.Parse(query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Parsed Filter: %s\n", parsed.Query[0].Filters.String())

	// Example 2: Using comparison operators in filters
	fmt.Println("\n=== Example 2: Comparison Operators in Filters ===")
	
	// Greater than example
	gtQuery := `sum:response_time{value:>500}`
	gtParsed, err := parser.Parse(gtQuery)
	if err != nil {
		panic(err)
	}
	
	// Less than example
	ltQuery := `sum:error_rate{value:<0.01}`
	ltParsed, err := parser.Parse(ltQuery)
	if err != nil {
		panic(err)
	}
	
	// Greater than or equal example
	gteQuery := `sum:cpu{usage:>=90}`
	gteParsed, err := parser.Parse(gteQuery)
	if err != nil {
		panic(err)
	}
	
	// Less than or equal example
	lteQuery := `sum:memory{usage:<=75}`
	lteParsed, err := parser.Parse(lteQuery)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Greater Than: %s -> %s\n", gtQuery, gtParsed.Query[0].Filters.String())
	fmt.Printf("Less Than: %s -> %s\n", ltQuery, ltParsed.Query[0].Filters.String())
	fmt.Printf("Greater Than or Equal: %s -> %s\n", gteQuery, gteParsed.Query[0].Filters.String())
	fmt.Printf("Less Than or Equal: %s -> %s\n", lteQuery, lteParsed.Query[0].Filters.String())

	// Example 3: Using regex filters
	fmt.Println("\n=== Example 3: Regex Filters ===")
	regexQuery := `sum:http.requests{path:~^/api/v[0-9]/users}`
	regexParsed, err := parser.Parse(regexQuery)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Regex Filter: %s -> %s\n", regexQuery, regexParsed.Query[0].Filters.String())

	// Example 4: Building filters programmatically
	fmt.Println("\n=== Example 4: Building Filters Programmatically ===")
	
	// Create a filter structure for: env:prod AND service:api
	filter := &ddqp.MetricFilter{}
	
	// First part: env:prod
	filter.Left = &ddqp.Param{}
	filter.Left.SimpleFilter = &ddqp.SimpleFilter{}
	filter.Left.SimpleFilter.FilterKey = "env"
	filter.Left.SimpleFilter.FilterSeparator = &ddqp.FilterSeparator{Colon: true}
	filter.Left.SimpleFilter.FilterValue = &ddqp.FilterValue{}
	filter.Left.SimpleFilter.FilterValue.SimpleValue = &ddqp.Value{}
	envVal := "prod"
	filter.Left.SimpleFilter.FilterValue.SimpleValue.Identifier = &envVal
	
	// Add AND separator
	filter.Parameters = append(filter.Parameters, &ddqp.Param{
		Separator: &ddqp.FilterValueSeparator{
			And: true,
		},
	})
	
	// Second part: service:api
	svcParam := &ddqp.Param{}
	svcParam.SimpleFilter = &ddqp.SimpleFilter{}
	svcParam.SimpleFilter.FilterKey = "service"
	svcParam.SimpleFilter.FilterSeparator = &ddqp.FilterSeparator{Colon: true}
	svcParam.SimpleFilter.FilterValue = &ddqp.FilterValue{}
	svcParam.SimpleFilter.FilterValue.SimpleValue = &ddqp.Value{}
	svcVal := "api"
	svcParam.SimpleFilter.FilterValue.SimpleValue.Identifier = &svcVal
	
	filter.Parameters = append(filter.Parameters, svcParam)
	
	// Create a query using this filter
	newQuery := &ddqp.MetricQuery{}
	newQuery.Query = []*ddqp.Query{{
		Aggregator: "avg",
		MetricName: "system.cpu.user",
		Filters:    filter,
		Grouping:   []string{"host"},
		By:         "by",
	}}
	
	fmt.Printf("Programmatically Built Query: %s\n", newQuery.String())
}