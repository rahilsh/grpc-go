// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	pb "rahilsh/grpc-go/proto"
	"time"
)

func main() {
	// Set up a connection to the server.
	var host string
	var hostHeader string
	var name string

	errorMessage := "Usage: /client -host <host> [-hostHeader] <hostHeader> [-name] <name>"

	flag.StringVar(&host, "host", "", errorMessage)
	flag.StringVar(&hostHeader, "hostHeader", "", errorMessage)
	flag.StringVar(&name, "name", "world", errorMessage)
	flag.Parse()

	if len(host) == 0 {
		fmt.Println(errorMessage)
		os.Exit(1)
	}

	println("host: " + host)
	println("hostHeader: " + hostHeader)
	println("name: " + name)

	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	if len(hostHeader) != 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "host", hostHeader)
	}
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not SayHello: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
