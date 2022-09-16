# DDQP (DataDog Query Parser)

This package is intended to make the parsing of DataDog queries easier by providing a parser around which clients/tooling can be built. Looking around the DataDog community it was difficult to find a package for programatically interacting with existing queries or building up new queries. Thus this package was developed in an attempt to make parsing queries easier.

This package is not a client, it just provides structure around which clients or other tools can be built.

## Architecture

This package is built around [`participle`](https://github.com/alecthomas/participle), which takes away the effort of building a full parser and instead allows us to focus on capturing the variations present in the DataDog query "language".

## Developmemt

This project uses [Hermit](https://cashapp.github.io/hermit/) for managing dependencies and go tools.

Each portion of the parser is divided into its own file, and each file has an associated test file for validating the parser discreet pieces. The end goal is to combine each unit of a query into one top-level struct that encapsulates a full query. For example, a `MetricQuery` parser is built upon a `MetricFilter` parser.

## Supported Queries

Currently there is only support for a simple metric query, but the goal is to layer in more functionality.
