package lesson_23_test

import (
	"context"
	"testing"
	"time"

	"github.com/artsadert/lesson_23/internal/interface/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestListNodesStreamWithTimeout(t *testing.T) {
	// Connect to the server (adjust address if needed)
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAuthServiceClient(conn)

	// 1. Login
	loginCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loginResp, err := client.Login(loginCtx, &proto.LoginRequest{
		Name:     "hello",
		Password: "test",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	accessToken := loginResp.AccessToken

	// 2. Prepare context with a short deadline (3 seconds)
	streamCtx, streamCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer streamCancel()

	// 3. Attach token to metadata
	md := metadata.Pairs("authorization", "Bearer "+accessToken)
	streamCtx = metadata.NewOutgoingContext(streamCtx, md)

	// 4. Call the streaming method
	stream, err := client.ListNodesStream(streamCtx, &proto.ListNodesRequest{})
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}

	// 5. Try to receive nodes – we expect a DeadlineExceeded error
	var received int
	for {
		_, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.DeadlineExceeded {
				t.Logf("Deadline exceeded as expected after %d nodes", received)
				return // test passes
			}
			t.Fatalf("Unexpected error: %v", err)
		}
		received++
	}
}
