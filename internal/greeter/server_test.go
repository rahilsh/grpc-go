package greeter

import (
	"context"
	"testing"

	pb "github.com/rahilsh/golang-lab/internal/greeterpb"
)

func TestSayHello(t *testing.T) {
	response, err := NewServer().SayHello(context.Background(), &pb.HelloRequest{Name: "Alice"})
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}
	if got, want := response.GetMessage(), "Hello Alice"; got != want {
		t.Errorf("SayHello() message = %q, want %q", got, want)
	}
}
