# Contributing to envx

Thank you for your interest in contributing to envx! This document provides guidelines for contributors.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/envx.git`
3. Create a feature branch: `git checkout -b feature-name`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit your changes: `git commit -am 'Add some feature'`
7. Push to your branch: `git push origin feature-name`
8. Create a Pull Request

## Development Setup

```bash
# Clone the repository
git clone https://github.com/Justblue0312/envx.git
cd envx

# Install dependencies
go mod download

# Run tests
go test ./...

# Run benchmarks
go test -bench=.

# Run linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

## Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed

## Testing

- All new features must include tests
- Maintain high test coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

## Submitting Changes

1. Ensure your code follows the style guidelines
2. All tests pass: `go test ./...`
3. No lint issues: `golangci-lint run`
4. Update documentation if needed
5. Submit a Pull Request with a clear description

## Release Process

Releases are automated through GitHub Actions:

1. Push a tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions will automatically:
   - Run all tests
   - Build binaries for multiple platforms
   - Create a GitHub release
   - Publish packages

## Questions?

Feel free to open an issue for questions or discussion.