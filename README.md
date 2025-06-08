# DDQP (DataDog Query Parser)

DDQP is a Go package for parsing and constructing DataDog queries programmatically. It provides a structured way to work with different types of DataDog queries without having to manually handle string manipulation or complex regex patterns.

[![Go Reference](https://pkg.go.dev/badge/github.com/jonwinton/ddqp.svg)](https://pkg.go.dev/github.com/jonwinton/ddqp)

## Features

- Parse DataDog queries into structured Go objects
- Generate queries programmatically with type safety
- Support for complex expressions, filters, and conditions
- Handle various query types (metrics, monitors)

## Installation

```bash
go get github.com/jonwinton/ddqp
```

## Usage

### Basic Metric Query Parsing

```go
package main

import (
    "fmt"
    "github.com/jonwinton/ddqp"
)

func main() {
    parser := ddqp.NewMetricQueryParser()
    query, err := parser.Parse("sum:system.cpu.user{host:web-* AND env:prod} by {host}")
    if err != nil {
        panic(err)
    }
    
    // Access structured data
    fmt.Printf("Aggregator: %s\n", query.Query[0].Aggregator)
    fmt.Printf("Metric Name: %s\n", query.Query[0].MetricName)
    
    // Convert back to string
    fmt.Printf("Query String: %s\n", query.String())
}
```

### Monitor Query Parsing

```go
parser := ddqp.NewMetricMonitorParser()
monitor, err := parser.Parse("avg(last_5m):system.cpu.user{env:prod} > 80")
if err != nil {
    panic(err)
}

fmt.Printf("Aggregation: %s\n", monitor.Aggregation)
fmt.Printf("Window: %s\n", monitor.EvaluationWindow)
fmt.Printf("Threshold: %g\n", monitor.Threshold)
fmt.Printf("Comparator: %s\n", monitor.Comparator)
```

### Complex Expressions

```go
parser := ddqp.NewMetricExpressionParser()
expression, err := parser.Parse("(sum:system.cpu.user{*} / sum:system.cpu.idle{*}) * 100")
if err != nil {
    panic(err)
}

// Formulas can be used to better understand expressions
formula := ddqp.NewMetricExpressionFormula(expression)
fmt.Printf("Formula: %s\n", formula.Formula)
for k, v := range formula.Expressions {
    fmt.Printf("%s = %s\n", k, v)
}
```

## Architecture

DDQP is built around [`participle`](https://github.com/alecthomas/participle), a parser library that makes it easy to define parsers from Go struct definitions. This allows DDQP to focus on capturing the variations present in the DataDog query language.

The parser is divided into separate components:

- **MetricQuery**: Basic DataDog metric queries
- **MetricFilter**: Filter expressions for queries (e.g., `{host:web-* AND env:prod}`)
- **MetricMonitor**: Monitor queries with evaluation windows and thresholds
- **MetricExpression**: Mathematical expressions involving metrics

## Development

This project uses [Hermit](https://cashapp.github.io/hermit/) for managing dependencies and Go tools.

```bash
# Activate Hermit environment
source ./bin/activate-hermit

# Run tests
./scripts/test

# Run linter
./bin/golangci-lint run
```

## Supported Features

- **Metric queries** with filtering and grouping
- **Monitor queries** with thresholds and comparisons
- **Mathematical expressions** combining multiple metrics
- **Complex filters** with AND/OR/NOT logic
- **Comparison operators** (>, <, >=, <=)
- **Regex filters** using `:~` operator

## Examples

See the [`_examples`](./_examples/) directory for more usage examples:

- [Basic Metrics](./examples/metrics/): Working with metric queries
- [Filters](./examples/filters/): Advanced filter usage
- [Monitors](./examples/monitors/): Creating and parsing monitors
- [Expressions](./examples/expressions/): Complex mathematical expressions

## License

This project is licensed under the terms found in the [LICENSE](./LICENSE) file.
