# idflake Example

This directory contains a small runnable program that demonstrates the common
ways to use `github.com/Shudharshan07/idflake/v2`.

Run it from this directory:

```sh
go run .
```

The example shows:

- creating a generator for a node
- generating one ID and a batch of IDs
- extracting timestamp, node, and sequence fields
- converting an ID to Base36 and parsing it back
- marshaling IDs safely through JSON
- using the SQL `Value` and `Scan` helpers

The local `replace` directive in `go.mod` points at the parent package, so this
example runs against the checkout you are editing.
