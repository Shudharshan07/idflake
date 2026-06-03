![Project Banner](.github/Idflake.png)

# idflake: Fast Distributed ID Generator for Go

idflake is a Go library for generating unique, time ordered int64 IDs.
You use it in backend services where ordering, speed, and safety matter.

Each ID encodes time, node, and sequence data.
Each call runs without locks.
Each ID works across JSON, databases, and text formats.

## Table of Contents

- [Main Features](#main-features)
- [Design](#design)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Batch Generation](#batch-generation)
- [Concurrency Model](#concurrency-model)
- [Encoding and Conversion](#encoding-and-conversion)
- [JSON Support](#json-support)
- [Database Support](#database-support)
- [Extracting Components](#extracting-components)
- [Limits](#limits)
- [Testing](#testing)
- [License](#license)

## Main Features

* Time ordered int64 IDs
* Lock free generation using atomics
* Monotonic clock support
* Configurable epoch and bit layout
* Batch ID generation
* Built in encoding helpers
* JSON safe representation
* SQL database integration

## Design

Each ID packs three fields into one int64 value.

* Timestamp in milliseconds since a custom epoch
* Node identifier for distributed systems
* Sequence number for collisions within the same millisecond

The generator relies on a monotonic time source.
Clock changes do not break ordering.

Default layout:

* Node bits: 10
* Sequence bits: 12
* Max nodes: 1024
* Max IDs per millisecond per node: 4096

## Getting Started

Create a generator with a node ID.
Call Generate to get an ID.

```go
sf := NewSnowFlake(1)
id := sf.Generate()
```

IDs increase over time on a single node.

## Configuration

You adjust layout using options.

Available options:

* WithEpoch
* WithNodeBits
* WithSequenceBits

Example:

```go
sf := NewSnowFlake(
    3,
    WithNodeBits(8),
    WithSequenceBits(14),
)
```

## Batch Generation

You generate many IDs in one call.
Returned IDs stay consecutive.

```go
ids, err := sf.GenerateN(100)
```

Batch size caps at the sequence capacity per millisecond.

## Concurrency Model

* Atomic state only
* No mutexes
* Safe across goroutines
* Strict ordering per node

A single generator instance is safe to share across all goroutines on a node.

## Encoding and Conversion

ID supports multiple formats.
Each format round trips without loss.

Supported formats:

* Decimal string
* Binary string
* Base32
* Base36
* Base64 from decimal bytes
* Hex
* Big endian 8 byte array

Examples:

```go
s := id.String()
id2, _ := ParseString(s)

b := id.IntBytes()
id3, _ := ParseIntBytes(b)
```

## JSON Support

IDs marshal as JSON strings.
This preserves full int64 precision.

```go
data, _ := json.Marshal(id)
json.Unmarshal(data, &id2)
```

## Database Support

ID implements sql.Valuer and sql.Scanner.
You store IDs as int64 columns.

```go
db.Exec("insert into table (id) values (?)", id)
```

## Extracting Components

You read fields back from an ID.

Available helpers:

* Timestamp
* Node
* Sequence

```go
ts := sf.Timestamp(id)
node := sf.Node(id)
seq := sf.Sequence(id)
```

## Limits

* IDs fit in signed int64
* IDs turn negative far in the future
* One generator maps to one node ID
* GenerateN caps at sequence capacity

## Testing

The test suite covers:

* Ordering
* Concurrency
* Batch generation
* All encodings
* JSON round trips
* Database integration
* Bit field reconstruction

Benchmarks show constant time generation under load.

## License

See LICENSE file for details.
