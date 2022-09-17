# dotodag-ql

**INSTALL HERMIT**

The `_examples` directory is a small place to build out little test examples, but there's also a Go tests option.

First time, go into `_examples/metrics` and run `go run main.go`. It'll send a simple query through the parser and print the AST.

Or checkout `metricquery_test.go` and run the tests with `./script/test`.


## Plan

- Build discreet units (example: a metric query parser)
- Build a higher order parser which allows for parsing **monitor queries** which just wraps the original metric query
- Build an even higher order parser which allows for formulas of metric queries and monitor queries
