// Example of using DDQP to work with DataDog metric expressions
package main

import (
	"fmt"
	"github.com/jonwinton/ddqp"
)

// This example demonstrates working with DataDog metric expressions:
// 1. Parsing complex mathematical expressions involving metrics
// 2. Converting expressions to formulas for better understanding
// 3. Working with nested expressions and parentheses
// 4. Creating expressions programmatically
func main() {
	// Create an expression parser
	parser := ddqp.NewMetricExpressionParser()

	// Example 1: Basic arithmetic operations
	fmt.Println("=== Example 1: Basic Arithmetic Operations ===")
	
	// Addition
	addExpr := `sum:system.cpu.user{*} + sum:system.cpu.system{*}`
	_, err := parser.Parse(addExpr)
	if err != nil {
		panic(err)
	}
	
	// Subtraction
	subExpr := `sum:system.cpu.user{*} - sum:system.cpu.idle{*}`
	_, err = parser.Parse(subExpr)
	if err != nil {
		panic(err)
	}
	
	// Multiplication
	mulExpr := `sum:system.memory.used{*} * 100`
	_, err = parser.Parse(mulExpr)
	if err != nil {
		panic(err)
	}
	
	// Division
	divExpr := `sum:system.network.bytes_sent{*} / 1024`
	_, err = parser.Parse(divExpr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Addition: %s\n", addExpr)
	fmt.Printf("Subtraction: %s\n", subExpr)
	fmt.Printf("Multiplication: %s\n", mulExpr)
	fmt.Printf("Division: %s\n", divExpr)

	// Example 2: Complex expressions with multiple operators
	fmt.Println("\n=== Example 2: Complex Expressions with Multiple Operators ===")
	complexExpr := `sum:metric.name{foo:bar} * sum:metric.name_two{foo:bar} + sum:metric.name_three{foo:bar} - 50`
	
	complexParsed, err := parser.Parse(complexExpr)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Complex Expression: %s\n", complexExpr)
	fmt.Printf("Parsed Expression: %s\n", complexParsed.String())

	// Example 3: Using parentheses for grouping operations
	fmt.Println("\n=== Example 3: Using Parentheses for Grouping Operations ===")
	parenExpr := `(sum:metric.name{foo:bar} - sum:metric.name_two{foo:bar}) / sum:metric.name_three{foo:bar}`
	
	parenParsed, err := parser.Parse(parenExpr)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Parenthesized Expression: %s\n", parenExpr)
	fmt.Printf("Parsed Expression: %s\n", parenParsed.String())

	// Example 4: Nested parentheses
	fmt.Println("\n=== Example 4: Nested Parentheses ===")
	nestedExpr := `(sum:metric.name{foo:bar} - (sum:metric.name_two{foo:bar} * 2)) / 10`
	
	nestedParsed, err := parser.Parse(nestedExpr)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Nested Expression: %s\n", nestedExpr)
	fmt.Printf("Parsed Expression: %s\n", nestedParsed.String())

	// Example 5: Converting expressions to formulas
	fmt.Println("\n=== Example 5: Converting Expressions to Formulas ===")
	formulaExpr := `(sum:system.cpu.user{*} / sum:system.cpu.idle{*}) * 100`
	
	formulaParsed, err := parser.Parse(formulaExpr)
	if err != nil {
		panic(err)
	}
	
	// Convert to formula
	formula := ddqp.NewMetricExpressionFormula(formulaParsed)
	
	fmt.Printf("Expression: %s\n", formulaExpr)
	fmt.Printf("Formula: %s\n", formula.Formula)
	fmt.Println("Expressions Map:")
	for k, v := range formula.Expressions {
		fmt.Printf("  %s = %s\n", k, v)
	}

	// Example 6: Percentage calculation
	fmt.Println("\n=== Example 6: Percentage Calculation ===")
	percentExpr := `(sum:system.disk.used{*} / sum:system.disk.total{*}) * 100`
	
	percentParsed, err := parser.Parse(percentExpr)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Percentage Expression: %s\n", percentExpr)
	fmt.Printf("Parsed Expression: %s\n", percentParsed.String())
}