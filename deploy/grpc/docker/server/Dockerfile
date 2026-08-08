# Build from the repository root:
#   docker build -t grpc-server:local -f deploy/grpc/docker/server/Dockerfile .

FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /grpc-server ./cmd/grpc-server

FROM scratch
COPY --from=builder /grpc-server /grpc-server
USER 65532:65532
EXPOSE 5050
ENTRYPOINT ["/grpc-server"]
