// Package main implements a client for the Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/rahilsh/golang-lab/internal/greeterpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	host := flag.String("host", "localhost:5050", "gRPC server address")
	hostHeader := flag.String("host-header", "", "optional host metadata value")
	name := flag.String("name", "world", "name to greet")
	flag.Parse()

	conn, err := grpc.NewClient(*host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if *hostHeader != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "host", *hostHeader)
	}

	response, err := pb.NewGreeterClient(conn).SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("say hello: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())
}
