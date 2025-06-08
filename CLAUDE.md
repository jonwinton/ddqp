# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DDQP (DataDog Query Parser) is a Go package for parsing DataDog queries. It provides a structured way to programmatically interact with and build DataDog queries. The package is not a client but provides a foundation upon which clients or other tools can be built.

## Architecture

The package is built around [participle](https://github.com/alecthomas/participle), which handles the parsing complexity, allowing the package to focus on capturing variations in the DataDog query language.

The parser is divided into separate files, each handling a specific part of the query structure:
- `metricexpression.go`: Handles metric expressions in queries
- `metricfilter.go`: Handles filtering metrics
- `metricmonitor.go`: Handles monitor-related functionality
- `metricquery.go`: Top-level struct that combines components into a complete query
- `parser.go`: Core parser functionality

Each component has associated test files for validation.

## Development Environment

This project uses [Hermit](https://cashapp.github.io/hermit/) for managing dependencies and Go tools. The binary tools are located in the `bin` directory and include:
- Go (1.21.3)
- golangci-lint (1.52.2)
- gotestsum (1.12.0)

## Common Commands

### Running Tests

```bash
# Run all tests with verbose output
./scripts/test

# Alternative way to run tests
./bin/gotestsum -f standard-verbose

# Run specific tests
./bin/go test -v ./path/to/package -run TestName
```

### Linting

```bash
# Run linter
./bin/golangci-lint run
```

### Building

```bash
# Build the package
./bin/go build

# Build example
./bin/go build -o ./bin/example ./_examples/metrics
```

### Working with Hermit

```bash
# Activate Hermit environment
source ./bin/activate-hermit
```

## Current Limitations

The package currently only supports simple metric queries. Future development aims to add more functionality and query types.