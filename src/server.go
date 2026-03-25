// Package main implements a server for Greeter service.
package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "rahilsh/grpc-go/proto"
)

const (
	port = ":5050"
)

// server is used to implement proto.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements proto.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		hostHeader := md.Get("host")
		if len(hostHeader) != 0 {
			println("host header: " + hostHeader[0])
		} else {
			println("host header not present!!!")
		}

	}
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
