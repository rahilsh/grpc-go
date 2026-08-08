// Package main implements a server for the Greeter service.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/rahilsh/grpc-go/internal/greeter"
	pb "github.com/rahilsh/grpc-go/proto"
	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", ":5050", "address on which to listen")
	flag.Parse()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, greeter.NewServer())
	log.Printf("server listening at %s", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
