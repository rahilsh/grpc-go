# Build from the repository root:
#   docker build -t grpc-client:local -f deploy/grpc/docker/client/Dockerfile .

FROM golang:1.27-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /grpc-client ./cmd/grpc-client

FROM scratch
COPY --from=builder /grpc-client /grpc-client
USER 65532:65532
ENTRYPOINT ["/grpc-client"]
