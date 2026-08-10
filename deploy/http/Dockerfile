# syntax=docker/dockerfile:1

# Build from the repository root:
#   docker build -t http-server:local -f deploy/http/Dockerfile .

# --- build stage ---
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/http-server ./cmd/http-server

# --- runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/http-server /http-server

EXPOSE 3000
ENV SERVER_ADDR=:3000
USER nonroot:nonroot
ENTRYPOINT ["/http-server"]
