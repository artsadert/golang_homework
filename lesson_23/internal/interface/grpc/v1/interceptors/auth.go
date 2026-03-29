package auth

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextKey string

const UserUUIDKey ContextKey = "user_uuid"

func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	return parts[1], nil
}

// ValidateAccessToken decodes and validates the token using the access secret.
// Returns the private claims (map) and an error.
func ValidateAccessToken(tokenStr string, config *entities.Config) (string, error) {
	token, err := config.JWTAccessSecret.Decode(tokenStr)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	// Check expiration
	if exp, is_exp := token.Expiration(); !is_exp || exp.Before(time.Now()) {
		return "", status.Error(codes.Unauthenticated, "token expired")
	}

	var userUUID string

	err = token.Get("user_uuid", &userUUID)
	if err != nil {
		log.Println(err)
		return "", status.Error(codes.Unauthenticated, "user_uuid not found in token")
	}

	return userUUID, nil
}

func UnaryAuthInterceptor(config *entities.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for Login and Refresh endpoints
		if info.FullMethod == "/proto.AuthService/Login" ||
			info.FullMethod == "/proto.AuthService/Refresh" {
			return handler(ctx, req)
		}

		tokenStr, err := extractTokenFromMetadata(ctx)
		if err != nil {
			return nil, err
		}

		userUUID, err := ValidateAccessToken(tokenStr, config)
		if err != nil {
			return nil, err
		}

		newCtx := context.WithValue(ctx, UserUUIDKey, userUUID)
		return handler(newCtx, req)
	}
}

func StreamAuthInterceptor(config *entities.Config) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info.FullMethod == "/proto.AuthService/Login" ||
			info.FullMethod == "/proto.AuthService/Refresh" {
			return handler(srv, ss)
		}

		ctx := ss.Context()

		tokenStr, err := extractTokenFromMetadata(ctx)
		if err != nil {
			return err
		}

		userUUID, err := ValidateAccessToken(tokenStr, config)
		if err != nil {
			return err
		}

		newCtx := context.WithValue(ctx, UserUUIDKey, userUUID)
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}
		return handler(srv, wrapped)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
