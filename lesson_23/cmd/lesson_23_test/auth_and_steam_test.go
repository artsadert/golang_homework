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

func TestAuthAndStreaming(t *testing.T) {
	// Connect to the gRPC server (adjust address if needed)
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAuthServiceClient(conn)

	// Obtain a valid token for tests that require authentication
	loginCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	loginResp, err := client.Login(loginCtx, &proto.LoginRequest{
		Name:     "hello",
		Password: "test",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	validToken := loginResp.AccessToken

	// Test: valid token + streaming with short timeout (should trigger DeadlineExceeded)
	t.Run("Valid token and timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		md := metadata.Pairs("authorization", "Bearer "+validToken)
		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, err := client.ListNodesStream(ctx, &proto.ListNodesRequest{})
		if err != nil {
			t.Fatalf("Failed to start stream: %v", err)
		}

		var nodesReceived int
		for {
			_, err := stream.Recv()
			if err != nil {
				if status.Code(err) == codes.DeadlineExceeded {
					t.Logf("Deadline exceeded after %d nodes (expected)", nodesReceived)
					break // test passes
				}
				t.Fatalf("Unexpected error: %v", err)
			}
			nodesReceived++
		}
	})

	// Test: invalid bearer token -> Unauthenticated
	t.Run("Invalid token", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		invalidToken := "this.is.an.invalid.token"
		md := metadata.Pairs("authorization", "Bearer "+invalidToken)
		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, err := client.ListNodesStream(ctx, &proto.ListNodesRequest{})
		if err != nil {
			t.Fatalf("Initial call error: %v", err)
		}
		_, err = stream.Recv()
		if err == nil {
			t.Fatal("Expected error on Recv, got nil")
		}
		if status.Code(err) != codes.Unauthenticated {
			t.Fatalf("Expected Unauthenticated, got %v", status.Code(err))
		}
	})

	// Test: no token at all -> Unauthenticated
	t.Run("No token", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// No metadata attached
		md := metadata.New(map[string]string{})
		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, err := client.ListNodesStream(ctx, &proto.ListNodesRequest{})
		if err != nil {
			t.Fatalf("Initial call error: %v", err)
		}
		_, err = stream.Recv()
		if err == nil {
			t.Fatal("Expected error on Recv, got nil")
		}
		if status.Code(err) != codes.Unauthenticated {
			t.Fatalf("Expected Unauthenticated, got %v", status.Code(err))
		}
	})

	// Test: context timeout on Login (to show deadline propagation)
	t.Run("Short timeout on login", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		_, err := client.Login(ctx, &proto.LoginRequest{
			Name:     "hello",
			Password: "test",
		})
		if err == nil {
			t.Fatal("Expected deadline error, got nil")
		}
		if status.Code(err) != codes.DeadlineExceeded {
			t.Fatalf("Expected DeadlineExceeded, got %v", status.Code(err))
		}
	})
}
