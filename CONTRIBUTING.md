# Contributing

Thank you for helping improve this example. Keep changes focused and open an
issue before starting a large redesign.

## Development

1. Install Go 1.25 or later.
2. Fork and clone the repository.
3. Create a branch from `main`.
4. Make the change and add or update tests for changed behavior.
5. Run the checks below before opening a pull request.

```sh
gofmt -w .
go mod tidy
go vet ./...
go test -race ./...
go build ./...
```

When changing `proto/hello.proto`, follow the generation instructions in the
[README](README.md) and commit the generated Go files. Do not edit generated
files directly.

## Pull Requests

- Explain the problem and the chosen solution.
- Link related issues.
- Keep unrelated changes out of the pull request.
- Update documentation when behavior or usage changes.
- Confirm that your contribution may be distributed under Apache-2.0.

By participating, you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).
