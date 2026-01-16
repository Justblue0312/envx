# envx

A powerful environment variable parsing library for Go applications with support for nested structs, custom decoders, and cross-platform compatibility.

## Features

- ✅ **Simple API** - Easy to use with struct tags
- ✅ **Nested Structs** - Support for nested configuration with double underscore separator
- ✅ **Custom Decoders** - Implement custom decoding logic
- ✅ **Type Support** - Strings, ints, bools, slices, maps, duration, and more
- ✅ **Cross-Platform** - Works on Windows, Linux, and macOS
- ✅ **Validation** - Required fields and type validation
- ✅ **Defaults** - Default values for missing environment variables
- ✅ **Prefix Support** - Group related configuration

## Installation

```bash
go get github.com/envx/envx
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/envx/envx"
)

type Config struct {
    Name    string `envconfig:"APP_NAME"`
    Port    int    `envconfig:"APP_PORT"`
    Debug   bool   `envconfig:"DEBUG" default:"false"`
    Timeout string `envconfig:"TIMEOUT" default:"30s"`
}

func main() {
    var config Config
    
    if err := envx.Process("", &config); err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Printf("App: %s running on port %d\n", config.Name, config.Port)
}
```

## Nested Structs

Use the `nested:"true"` tag to enable double underscore separation for nested structs:

```go
type DatabaseConfig struct {
    Host     string `envconfig:"HOST" nested:"true"`
    Port     int    `envconfig:"PORT" nested:"true"`
    Password string `envconfig:"PASSWORD" nested:"true"`
}

type Config struct {
    AppName  string         `envconfig:"APP_NAME"`
    Database DatabaseConfig `nested:"true"`
}
```

Environment variables:
- `APP_NAME=myapp`
- `DATABASE__HOST=localhost`
- `DATABASE__PORT=5432`
- `DATABASE__PASSWORD=secret`

## Supported Types

- **Basic Types**: `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `bool`, `float32`, `float64`
- **Complex Types**: `[]string`, `[]int`, `map[string]string`, `time.Duration`
- **Pointers**: Pointer to any supported type
- **Custom Types**: Types implementing `Decoder` or `Setter` interfaces

## Custom Decoders

Implement the `Decoder` interface for custom parsing logic:

```go
type CustomType struct {
    Value string
}

func (c *CustomType) Decode(value string) error {
    c.Value = "decoded:" + value
    return nil
}
```

## Struct Tags

- `envconfig:"VAR_NAME"` - Custom environment variable name
- `default:"value"` - Default value if environment variable is not set
- `required:"true"` - Mark field as required (error if not set)
- `nested:"true"` - Enable nested struct with double underscore separator
- `ignored:"true"` - Skip field during processing
- `split_words:"true"` - Convert CamelCase to SNAKE_CASE automatically

## Cross-Platform Support

envx automatically handles platform differences:

- **Windows**: Case-insensitive environment variable lookup
- **Unix/Linux**: Standard case-sensitive lookup

## API Reference

### `Process(prefix string, spec interface{}) error`

Populates the specified struct with environment variables.

- `prefix`: Optional prefix for environment variables
- `spec`: Pointer to struct to populate

### `MustProcess(prefix string, spec interface{})`

Same as `Process` but panics on error.

### `CheckDisallowed(prefix string, spec interface{}) error`

Checks for unknown environment variables with the given prefix.

## Examples

See the [examples](examples/) directory for complete working examples:

- [Basic Usage](examples/basic_usage/) - Simple configuration with prefix
- [Nested Structs](examples/nested_struct/) - Nested configuration example
- [Advanced Usage](examples/advanced_usage/) - Custom types and complex fields

## Performance

Benchmark results on modern hardware:

```
BenchmarkProcess-4   	  304969	      3996 ns/op
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.