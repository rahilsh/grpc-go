// Package main implements a server for the PingService.
package main

import (
	"flag"
	"log"
	"net"

	pb "github.com/rahilsh/golang-lab/internal/pingpb"
	"github.com/rahilsh/golang-lab/internal/pingserver"
	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", ":8080", "address on which to listen")
	flag.Parse()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterPingServiceServer(server, pingserver.NewServer())
	log.Printf("ping server listening at %s", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
