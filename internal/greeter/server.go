// Package greeter implements the Greeter service.
package greeter

import (
	"context"
	"log"

	pb "github.com/rahilsh/golang-lab/internal/greeterpb"
	"google.golang.org/grpc/metadata"
)

// Server implements the generated Greeter service interface.
type Server struct {
	pb.UnimplementedGreeterServer
}

// NewServer creates a Greeter service.
func NewServer() *Server {
	return &Server{}
}

// SayHello returns a greeting for the requested name.
func (s *Server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	if values := metadata.ValueFromIncomingContext(ctx, "host"); len(values) > 0 {
		log.Printf("host metadata: %s", values[0])
	}
	log.Printf("received: %s", request.GetName())
	return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}
