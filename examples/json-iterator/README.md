# JSON Iterator Example

This example demonstrates how to integrate [json-iterator/go](https://github.com/json-iterator/go) with Gin to replace the default encoding/json for better performance.

## How it works

Gin supports custom JSON serialization and deserialization logic via the `json.API` variable. By implementing the `json.Core` interface (which includes `Marshal`, `Unmarshal`, `NewEncoder`, `NewDecoder` etc.), we can swap out the underlying JSON engine.

## Usage

1. Define your custom configuration using `jsoniter.Config`.
2. Implement the `json.Core` interface wrappers.
3. Assign your custom implementation to `json.API` before creating the Gin engine.

## Run the example

```bash
go run main.go
```

Test it:

```bash
curl http://localhost:8080/ping
# Output: {"message":"pong"}
```
