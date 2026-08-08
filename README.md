# gRPC Go Hello World

[![Go](https://github.com/rahilsh/grpc-go/actions/workflows/go.yml/badge.svg)](https://github.com/rahilsh/grpc-go/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A small gRPC client and server written in Go, with Docker and Kubernetes examples.

> [!NOTE]
> This is an unofficial learning project, not the official
> [`grpc/grpc-go`](https://github.com/grpc/grpc-go) repository. It uses plaintext
> transport for local demonstration and is not a production-ready service.

## Requirements

- Go 1.25 or later
- Docker (optional)
- `kubectl` and a local cluster such as [kind](https://kind.sigs.k8s.io/) (optional)
- `protoc` and its Go plugins, only when changing `proto/hello.proto`

## Run Locally

Start the server:

```sh
go run ./cmd/server
```

In another terminal, call it with the client:

```sh
go run ./cmd/client -name Alice
```

The client prints `Greeting: Hello Alice`. Use `-host` to target another server
or `-host-header` to attach host metadata.

## Check Changes

Run the same core checks used by CI:

```sh
gofmt -w .
go vet ./...
go test -race ./...
go build ./...
```

## Generate Protobuf Code

Install [Protocol Buffers](https://protobuf.dev/installation/) and the pinned Go
plugins:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1
export PATH="$PATH:$(go env GOPATH)/bin"
```

After changing the schema, regenerate and commit both generated files:

```sh
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/hello.proto
```

Do not edit `proto/*.pb.go` directly.

## Run With Docker

Build both images:

```sh
docker build -t grpc-go-server:local -f docker/server/Dockerfile .
docker build -t grpc-go-client:local -f docker/client/Dockerfile .
```

Create a network, start the server, and run the client:

```sh
docker network create grpc-go
docker run --rm -d --name grpc-go-server --network grpc-go grpc-go-server:local
docker run --rm --network grpc-go grpc-go-client:local \
  -host grpc-go-server:5050 -name Alice
docker stop grpc-go-server
docker network rm grpc-go
```

## Run With kind

Build the images, load them into a running kind cluster, and apply the manifests:

```sh
kind load docker-image grpc-go-server:local grpc-go-client:local
kubectl apply -f kubernetes/server.yaml -f kubernetes/server-svc.yaml
kubectl apply -f kubernetes/client.yaml
kubectl logs job/grpc-go-client
```

Remove the example resources with `kubectl delete -f kubernetes/`.

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening
a pull request. For help, see [SUPPORT.md](SUPPORT.md). Report security concerns
according to [SECURITY.md](SECURITY.md).

## License

Licensed under the [Apache License 2.0](LICENSE). The protobuf example retains
the applicable copyright and license notice from the gRPC authors.
