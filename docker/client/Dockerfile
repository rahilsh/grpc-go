FROM golang:1.16 as builder

# Retrieve the dependencies.
RUN go get google.golang.org/grpc

# Copy local code to the container image.
WORKDIR /workdir
COPY . ./

# Build the command inside the container.
# To facilitate gRPC testing, this container includes a client command.

RUN CGO_ENABLED=0 go build -tags=grpc-go -o ./client src/client.go

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian

# Copy the binaries to the production image from the builder stage.
COPY --from=builder /workdir/client /client

CMD ["sleep", "infinity"]