.PHONY: build run-http run-http-client run-grpc run-grpc-client test test-cover lint vet fmt tidy proto clean docker

BINDIR := bin

build: ## Build all binaries
	go build -o $(BINDIR)/http-server ./cmd/http-server
	go build -o $(BINDIR)/http-client ./cmd/http-client
	go build -o $(BINDIR)/grpc-server ./cmd/grpc-server
	go build -o $(BINDIR)/grpc-client ./cmd/grpc-client

run-http: ## Run the HTTP server
	go run ./cmd/http-server

run-http-client: ## Run the HTTP client
	go run ./cmd/http-client

run-grpc: ## Run the gRPC server
	go run ./cmd/grpc-server

run-grpc-client: ## Run the gRPC client
	go run ./cmd/grpc-client -name Alice

test: ## Run tests
	go test -race ./...

test-cover: ## Run tests with coverage report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint: ## Run golangci-lint
	golangci-lint run

vet: ## Run go vet
	go vet ./...

fmt: ## Format code
	gofmt -w .

tidy: ## Tidy go.mod
	go mod tidy

proto: ## Regenerate protobuf code from proto/hello.proto
	protoc --go_out=. --go_opt=module=github.com/rahilsh/golang-lab \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/rahilsh/golang-lab \
	       proto/hello.proto

clean: ## Remove build artifacts
	rm -rf $(BINDIR) coverage.out

docker: ## Build all docker images (repository root as context)
	docker build -t http-server:local -f deploy/http/Dockerfile .
	docker build -t grpc-server:local -f deploy/grpc/docker/server/Dockerfile .
	docker build -t grpc-client:local -f deploy/grpc/docker/client/Dockerfile .
