# golang-lab

[![CI](https://github.com/rahilsh/golang-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/rahilsh/golang-lab/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A single Go module (`github.com/rahilsh/golang-lab`, requires Go 1.25+) with a
few small, self-contained example programs. The layout follows the official
[Organizing a Go module](https://go.dev/doc/modules/layout) guidance: all
binaries live in `cmd/`, private packages in `internal/`, and non-Go assets in
their own top-level directories.

## Layout

```text
cmd/
  http-server/   # net/http server
  http-client/   # HTTP client
  grpc-server/   # gRPC greeter server
  grpc-client/   # gRPC greeter client
  ping-server/   # streaming gRPC ping server
  ping-client/   # streaming gRPC ping client
  basics/        # language basics examples (-demo)
  functions/     # functions examples (-demo)
  types/         # types & interfaces examples (-demo)
  errors/        # errors & panic examples (-demo)
  concurrency/   # goroutines & channels examples (-demo)
  config/        # TOML configuration example
  web/           # JSON/HTTP examples (-demo)
  database/      # GORM, database/sql, and TiDB utilities (-demo)
internal/
  httpserver/    # HTTP routes, handlers, server setup (tested)
  greeter/       # gRPC Greeter service implementation (tested)
  greeterpb/     # generated protobuf/gRPC code
  pingserver/    # streaming PingService implementation
  pingpb/        # generated protobuf/gRPC code
  sqrt/          # simple, table-driven, and benchmark test examples (tested)
  demo/          # tiny -demo dispatcher shared by the example binaries
proto/
  hello.proto    # gRPC greeter schema
  ping.proto     # gRPC streaming schema
deploy/
  http/          # HTTP Dockerfile
  grpc/          # gRPC Dockerfiles and Kubernetes manifests
```

Install any command directly:

```sh
go install github.com/rahilsh/golang-lab/cmd/http-server@latest
go install github.com/rahilsh/golang-lab/cmd/grpc-server@latest
```

## HTTP example

```sh
go run ./cmd/http-server                     # listens on :3000 (SERVER_ADDR)
curl -i http://localhost:3000/               # 200: Hello world!
curl -i http://localhost:3000/books          # 200: {"total":0,"count":0,"books":[]}
go run ./cmd/http-client                      # GET CLIENT_URL (default jsonplaceholder)
```

| Variable      | Default                                        | Used by      |
| ------------- | ---------------------------------------------- | ------------ |
| `SERVER_ADDR` | `:3000`                                        | http-server  |
| `CLIENT_URL`  | `https://jsonplaceholder.typicode.com/todos/1` | http-client  |

## gRPC example

```sh
go run ./cmd/grpc-server                       # listens on :5050
go run ./cmd/grpc-client -name Alice           # prints: Greeting: Hello Alice
```

> [!NOTE]
> The gRPC example is an unofficial learning project, not the official
> [`grpc/grpc-go`](https://github.com/grpc/grpc-go) repository. It uses plaintext
> transport for local demonstration and is not production-ready.

Regenerate protobuf code after editing `proto/hello.proto` or `proto/ping.proto`
(do not edit the generated files in `internal/greeterpb` or `internal/pingpb`):

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1
export PATH="$PATH:$(go env GOPATH)/bin"
make proto
```

## Streaming gRPC example

`PingService` demonstrates a bidirectional streaming RPC (`PingStream`) alongside
a unary `Ping`:

```sh
go run ./cmd/ping-server                       # listens on :8080
go run ./cmd/ping-client -msg ping -count 3    # unary + streamed pongs
```

## Learning examples

Snippets collected while learning Go are grouped into a few runnable command
binaries by topic. Each bundles several demos selected with a `-demo` flag; run
without the flag to list what's available:

```sh
go run ./cmd/basics -demo fizzbuzz
go run ./cmd/concurrency -demo select
go run ./cmd/basics                    # lists available demos
```

| Command | Demos (`-demo`) |
| ------- | --------------- |
| `basics` | `hello-world`, `mean`, `if`, `switch`, `loops`, `fizzbuzz`, `strings`, `sprintf`, `even-ended-numbers`, `slices`, `slice-max`, `maps`, `word-count` |
| `functions` | `functions`, `function-parameters`, `returning-errors`, `defer`, `content-type` |
| `types` | `structs`, `receivers`, `methods`, `constructor`, `embedded-structs`, `interfaces`, `io-writer` |
| `errors` | `custom-errors`, `panic-recover`, `error-wrapping` |
| `concurrency` | `goroutines`, `channels`, `channel-content-type`, `select`, `md5-concurrent` |
| `web` | `json`, `http-get`, `github-api`, `httpd`, `kv-store` |
| `database` | `orm` (GORM), `dao` (`database/sql`), `region-splitter` (TiDB PD) |

`cmd/config` (single program) reads a TOML file. Testing patterns — simple,
table-driven (inline and CSV-backed), and benchmark — live together in
`internal/sqrt` and run via `go test ./internal/sqrt`.

### Resources

- LinkedIn Learning — Learning Go (course exercise files):
  <https://github.com/LinkedInLearning/learning-go-2875237>
- Go design patterns:
  <https://github.com/AgarwalConsulting/Go-Training/tree/master/patterns/design>
  ([slides](https://go-design-patterns.slides.algogrit.com/))
- [Effective Go](https://golang.org/doc/effective_go) ·
  [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Dave Cheney — [SOLID Go Design](https://dave.cheney.net/2016/08/20/solid-go-design)

## Checks

Run repository-wide from the root:

```sh
make fmt vet test build lint
# or directly:
gofmt -w . && go vet ./... && go test -race -cover ./... && go build ./... && golangci-lint run
```

## Docker

Images build from the repository root so they share the root `go.mod`/`go.sum`:

```sh
make docker
# or individually:
docker build -t http-server:local -f deploy/http/Dockerfile .
docker build -t grpc-server:local -f deploy/grpc/docker/server/Dockerfile .
docker build -t grpc-client:local -f deploy/grpc/docker/client/Dockerfile .
```

## Kubernetes (gRPC, via kind)

```sh
kind load docker-image grpc-server:local grpc-client:local
kubectl apply -f deploy/grpc/kubernetes/server.yaml -f deploy/grpc/kubernetes/server-svc.yaml
kubectl apply -f deploy/grpc/kubernetes/client.yaml
kubectl logs job/grpc-go-client
```

Remove the resources with `kubectl delete -f deploy/grpc/kubernetes/`.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md), [SUPPORT](SUPPORT.md), and
[SECURITY](SECURITY.md). Participation is governed by the
[Code of Conduct](CODE_OF_CONDUCT.md).

## License

Licensed under the [Apache License 2.0](LICENSE). The HTTP example originated
from the author's MIT-licensed `go-server` project and is redistributed here
under Apache-2.0. The gRPC protobuf files retain their upstream gRPC
Apache-2.0 notices.
