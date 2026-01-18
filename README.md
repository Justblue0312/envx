# envx

A powerful environment variable parsing library for Go applications with support for nested structs, custom decoders, and cross-platform compatibility.

## Features

- ✅ **Simple API** - Easy to use with struct tags
- ✅ **Nested Structs** - Support for nested configuration with single underscore separator
- ✅ **Custom Decoders** - Implement custom decoding logic
- ✅ **Type Support** - Strings, ints, bools, slices, maps, duration, and more
- ✅ **Cross-Platform** - Works on Windows, Linux, and macOS
- ✅ **Validation** - Required fields and type validation
- ✅ **Defaults** - Default values for missing environment variables
- ✅ **Prefix Support** - Group related configuration

## Installation

```bash
go get github.com/justblue0312/envx@latest
```

## Requirements

- Go 1.21 or later

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/justblue0312/envx"
)

type Config struct {
    Name    string `envx:"APP_NAME"`
    Port    int    `envx:"APP_PORT"`
    Debug   bool   `envx:"DEBUG" default:"false"`
    Timeout string `envx:"TIMEOUT" default:"30s"`
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

Use the `nested:"true"` tag to enable single underscore separation for nested structs. Only the main nested struct field needs the `nested:"true"` tag - fields within the nested struct are automatically detected:

```go
type DatabaseConfig struct {
    Host     string `envx:"HOST"`        // No nested tag needed
    Port     int    `envx:"PORT"`        // No nested tag needed
    Password string `envx:"PASSWORD"`    // No nested tag needed
}

type Config struct {
    AppName  string         `envx:"APP_NAME"`
    Database DatabaseConfig `nested:"true"`  // Only main node needs nested tag
}
```

Environment variables:
- `APP_NAME=myapp`
- `DATABASE_HOST=localhost`
- `DATABASE_PORT=5432`
- `DATABASE_PASSWORD=secret`

### Prefix Inheritance

When the main nested struct has an `envx` tag, nested fields automatically inherit the prefix:

```go
type DatabaseConfig struct {
    Host     string `envx:"HOST"`
    Port     int    `envx:"PORT"`
    Password string `envx:"PASSWORD"`
}

type Config struct {
    AppName  string         `envx:"APP_NAME"`
    Database DatabaseConfig `envx:"DB" nested:"true"`  // Has envx tag for prefix
}
```

Environment variables:
- `APP_NAME=myapp`
- `DB_HOST=localhost`         // Inherits DB prefix
- `DB_PORT=5432`              // Inherits DB prefix
- `DB_PASSWORD=secret`        // Inherits DB prefix

## Advanced Type Examples

### Time and URL Types

```go
type AdvancedConfig struct {
    Duration     time.Duration   `envx:"TIMEOUT"`
    Location     *time.Location `envx:"TIMEZONE"`
    DatabaseURL  *url.URL       `envx:"DATABASE_URL"`
}
```

Environment variables:
- `TIMEOUT=30s`
- `TIMEZONE=UTC`
- `DATABASE_URL=postgres://user:pass@localhost:5432/db`

### Slices and Maps

```go
type CollectionConfig struct {
    Tags         []string          `envx:"TAGS"`              // Comma-separated: "tag1,tag2,tag3"
    Numbers      []int             `envx:"NUMBERS"`           // Comma-separated: "1,2,3"
    Data         []byte            `envx:"DATA"`              // Raw bytes
    StringPtrs   []*string        `envx:"STRING_PTRS"`       // Comma-separated string pointers
    Properties   map[string]string `envx:"PROPERTIES"`        // "key1:val1,key2:val2"
    Scores       map[string]int    `envx:"SCORES"`            // "alice:100,bob:200"
}
```

### Custom Types with Interfaces

```go
// Using encoding.TextUnmarshaler
type Email string
func (e *Email) UnmarshalText(text []byte) error {
    *e = Email(strings.ToLower(string(text)))
    return nil
}

// Using Decoder interface
type DatabaseConfig struct {
    ConnectionString string
    MaxConnections int
}

func (d *DatabaseConfig) Decode(value string) error {
    // Custom parsing logic for database config
    parts := strings.Split(value, ";")
    if len(parts) != 2 {
        return fmt.Errorf("invalid database config format")
    }
    d.ConnectionString = parts[0]
    if max, err := strconv.Atoi(parts[1]); err == nil {
        d.MaxConnections = max
    }
    return nil
}
```

## Supported Types

- **Basic Types**: `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `bool`, `float32`, `float64`
- **Time Types**: `time.Duration`, `*time.Location`
- **URL Types**: `*url.URL`
- **Complex Types**: `[]string`, `[]int`, `[]byte`, slices of any supported types, slices of pointers to supported types
- **Maps**: `map[string]string`, `map[string]int`, and maps with any supported key/value types
- **Pointers**: Pointer to any supported type
- **Custom Types**: Types implementing `Decoder` or `Setter` interfaces
- **Interface Support**: Types implementing `encoding.TextUnmarshaler` and `encoding.BinaryUnmarshaler`

## Custom Types

### Decoder Interface

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

### Setter Interface

The `Setter` interface provides an alternative to `Decoder` for simpler cases:

```go
type CustomSetter struct {
    Data string
}

func (c *CustomSetter) Set(value string) error {
    c.Data = "set:" + value
    return nil
}
```

### When to Use Setter vs Decoder

**Use `Setter` when:**
- You only need to set a value from a string
- The logic is simple and straightforward
- You don't need access to the full field context

**Use `Decoder` when:**
- You need more complex parsing logic
- You need to validate the input value
- You might need to handle different input formats
- You want to return detailed error information

Both interfaces are checked, with `Decoder` taking precedence over `Setter`.

## Struct Tags

- `envx:"VAR_NAME"` - Custom environment variable name
- `default:"value"` - Default value if environment variable is not set
- `required:"true"` - Mark field as required (error if not set)
- `nested:"true"` - Enable nested struct with single underscore separator
- `ignored:"true"` - Skip field during processing
- `split_words:"true"` - Convert CamelCase to SNAKE_CASE automatically

## Cross-Platform Support

envx automatically handles platform differences:

- **Windows**: Case-insensitive environment variable lookup
- **Unix/Linux**: Standard case-sensitive lookup

## API Reference

##### `Process(prefix string, spec any) error`

Populates the specified struct with environment variables.

- `prefix`: Optional prefix for environment variables
- `spec`: Pointer to struct to populate

##### `MustProcess(prefix string, spec any)`

Same as `Process` but panics on error.

#### `CheckDisallowed(prefix string, spec any) error`

Checks for unknown environment variables with the given prefix.

## Examples

See the [examples](examples/) directory for complete working examples:

- [Basic Usage](examples/basic_usage/) - Simple configuration with prefix.
- [Nested Structs](examples/nested_struct/) - Nested configuration example.
- [Advanced Usage](examples/advanced_usage/) - Custom types and complex fields.
- [Comprehensive](examples/comprehensive_types/) - Comprehensive types.

## Release Process

envx follows semantic versioning and uses automated releases:

### Versioning

- **Major (X.0.0)**: Breaking changes
- **Minor (0.Y.0)**: New features (backward compatible)
- **Patch (0.0.Z)**: Bug fixes (backward compatible)

### Creating a Release

```bash
# Update version in go.mod (if needed)
# Commit your changes
git add .
git commit -m "Prepare release v1.0.0"

# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions will automatically:
- Run all tests
- Build binaries for multiple platforms
- Create a GitHub release
- Generate checksums

### Development Versions

For development builds:

```bash
make release
```

This creates snapshot binaries in the `dist/` directory.

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

For security issues, please see [SECURITY.md](SECURITY.md) for reporting procedures.
