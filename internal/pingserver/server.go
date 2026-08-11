// Package pingserver implements the PingService, demonstrating unary and
// bidirectional streaming gRPC handlers.
package pingserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/rahilsh/golang-lab/internal/pingpb"
	"google.golang.org/grpc"
)

// Server implements the generated PingService interface.
type Server struct {
	pb.UnimplementedPingServiceServer
}

// NewServer creates a PingService.
func NewServer() *Server {
	return &Server{}
}

// Ping returns a single pong for the request (unary RPC).
func (s *Server) Ping(_ context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("ping: %s", request.GetMsg())
	return &pb.Response{Msg: fmt.Sprintf("%s - pong", request.GetMsg())}, nil
}

// PingStream pongs every request until the client closes the stream
// (bidirectional streaming RPC).
func (s *Server) PingStream(stream grpc.BidiStreamingServer[pb.Request, pb.Response]) error {
	for {
		request, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Print("client closed the stream")
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("ping stream: %s", request.GetMsg())
		if err := stream.Send(&pb.Response{
			Msg: fmt.Sprintf("%s - pong at %s", request.GetMsg(), time.Now().Format(time.RFC3339)),
		}); err != nil {
			return err
		}
	}
}
