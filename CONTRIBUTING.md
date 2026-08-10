# Contributing

Thank you for helping improve these examples. Keep changes focused and open an
issue before starting a large redesign.

This is a single Go module (`github.com/rahilsh/golang-lab`). Commands live in
`cmd/`, private packages in `internal/`. Run tooling from the repository root so
it covers everything at once.

## Development

1. Install Go 1.25 or later.
2. Fork and clone the repository.
3. Create a branch from `main`.
4. Make the change and add or update tests for changed behavior.
5. Run the checks below from the repository root before opening a pull request.

```sh
gofmt -w .
go mod tidy
go vet ./...
go test -race ./...
go build ./...
golangci-lint run
```

When changing `proto/hello.proto`, regenerate with `make proto` and commit the
updated files in `internal/greeterpb`. Do not edit generated files directly.

## Pull Requests

- Explain the problem and the chosen solution.
- Link related issues.
- Keep unrelated changes out of the pull request.
- Update documentation when behavior or usage changes.
- Confirm that your contribution may be distributed under the repository's
  Apache-2.0 license.

By participating, you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).
