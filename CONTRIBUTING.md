# Contributing to poc-requests-go

## Code of Conduct

By participating in this project, you are expected to uphold our code of conduct. Please be respectful and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Git
- Make (optional, for using Makefile commands)

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make bench
```

## Release Process

### Creating a Release

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a git tag: `git tag -a v1.0.0 -m "Release version 1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. Create GitHub release with release notes

## Documentation

### Code Documentation

- Add godoc comments for all public functions and types
- Include examples in documentation where helpful
- Keep documentation up-to-date with code changes

### API Documentation

- Document all supported endpoints in README.md
- Include examples of usage
- Document authentication requirements
- Explain error handling
