// Package main implements a client for the PingService that exercises both the
// unary Ping and the bidirectional streaming PingStream RPCs.
package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/rahilsh/golang-lab/internal/pingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := flag.String("host", "localhost:8080", "gRPC server address")
	msg := flag.String("msg", "ping", "message to send")
	count := flag.Int("count", 3, "number of streamed messages")
	flag.Parse()

	conn, err := grpc.NewClient(*host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := pb.NewPingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unary RPC.
	response, err := client.Ping(ctx, &pb.Request{Msg: *msg})
	if err != nil {
		log.Fatalf("ping: %v", err)
	}
	log.Printf("unary reply: %s", response.GetMsg())

	// Bidirectional streaming RPC.
	stream, err := client.PingStream(ctx)
	if err != nil {
		log.Fatalf("ping stream: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			reply, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				log.Printf("recv: %v", err)
				return
			}
			log.Printf("stream reply: %s", reply.GetMsg())
		}
	}()

	for i := 0; i < *count; i++ {
		if err := stream.Send(&pb.Request{Msg: *msg}); err != nil {
			log.Fatalf("send: %v", err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("close send: %v", err)
	}
	<-done
}
