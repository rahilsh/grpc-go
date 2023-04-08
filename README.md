# grpc-go
Sample grpc server and client using golang

Quick start - https://grpc.io/docs/languages/go/quickstart/

### Install protobuf

#### Linux

```shell
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+
```

#### Mac
```shell
brew install protobuf
```

### Install the protocol compiler plugins for Go

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

### Update your PATH so that the protoc compiler can find the plugins
```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Generate code
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto
```

## Test locally

Run server
```shell
go run src/server.go
```

On another terminal make gRPC call using client
```shell
go run src/client.go -host localhost:5050
```

### Build images 
Build client image
```
docker build --tag 10.24.9.10/dsql/rahil/grpc-go-client:v0.0.5 --file docker/ClientDockerfile .
```

Build server image
```
docker build --tag 10.24.9.10/dsql/rahil/grpc-go-server:v0.0.4 --file docker/ServerDockerFile .
```

### Run on docker 

Create a network
```shell
docker network create grpc
```

Run client container 
```shell
docker run -d --name grpc-client --network=grpc 10.24.9.10/dsql/rahil/grpc-go-client:v0.0.6
```

Run server container
```shell
docker run -d --name grpc-server --network=grpc 10.24.9.10/dsql/rahil/grpc-go-server:v0.0.5
```

Get Server IP
```shell
docker inspect grpc-server | grep "IPAddress"
```

Exec into client container
```shell
docker exec -it 02b354b93a0819560bccca275d0017eca985525d02921158acd50cb46df04464 sh
```

Call Server
```shell
./client -host 172.20.0.2:5050
```

### Run on kubernetes

Push client image to harbor
```
docker push 10.24.9.10/dsql/rahil/grpc-go-client:v0.0.5
```

Push server image to harbor
```
docker push 10.24.9.10/dsql/rahil/grpc-go-server:v0.0.4
```

Create deployment, replicatset, services and pods on kubernetes
```shell
kubectl apply -f kubernetes
```

Exec to client pod
```shell
kubectl exec -n namespace pod/grpc-go-client-68bd85c465-bm6mx --stdin --tty -c grpc-go-client -- sh
```

Call server
```
/client -host grpc-go-server-svc:5050
```
