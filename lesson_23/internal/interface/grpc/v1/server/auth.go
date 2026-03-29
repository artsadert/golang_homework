package server

import (
	"context"
	"time"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/interface/grpc/proto"

	grpc_middleware "github.com/artsadert/lesson_23/internal/interface/grpc/v1/interceptors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	config      *entities.Config
	userService interfaces.UserService
}

func NewAuthServer(config *entities.Config, userService interfaces.UserService) *AuthServer {
	return &AuthServer{
		config:      config,
		userService: userService,
	}
}

// Login authenticates the user and returns tokens.
func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Call the existing user service to authenticate
	userResult, err := s.userService.Authenticate(&command.LoginUserCommand{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	// Generate tokens – reuse your existing helper (export it, or copy the logic)
	tokens, err := s.createTokens(userResult.Result.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &proto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.JWT_TTL.Seconds()),
	}, nil
}

// Refresh issues a new access token using a valid refresh token.
func (s *AuthServer) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.LoginResponse, error) {
	// Decode the refresh token using the refresh secret
	token, err := s.config.JWTRefreshSecret.Decode(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	var token_userUUID string

	err = token.Get("user_uuid", &token_userUUID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user_uuid not found in token")
	}

	userUUID, err := uuid.Parse(token_userUUID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user_uuid not valid in token")
	}

	// Optionally verify the user still exists (you might want to call userService.GetUser)
	_, err = s.userService.GetUser(userUUID)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	// Generate new tokens (access + refresh) – or just a new access token
	tokens, err := s.createTokens(userUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &proto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.JWT_TTL.Seconds()),
	}, nil
}

// Validate checks the validity of an access token.
func (s *AuthServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	userUUID, err := grpc_middleware.ValidateAccessToken(req.AccessToken, s.config)
	if err != nil {
		return &proto.ValidateResponse{Valid: false, Message: err.Error()}, nil
	}

	return &proto.ValidateResponse{Valid: true, UserUuid: userUUID}, nil
}

// createTokens is a helper that generates access and refresh tokens.
// (You can replace this with a call to an exported function from your existing jwt_helper.)
func (s *AuthServer) createTokens(userUUID uuid.UUID) (*struct {
	AccessToken  string
	RefreshToken string
}, error,
) {
	// Access token claims
	accessClaims := map[string]interface{}{
		"type":      "access",
		"exp":       time.Now().Add(s.config.JWT_TTL).Unix(),
		"iat":       time.Now().Unix(),
		"user_uuid": userUUID.String(),
	}
	_, accessTokenString, err := s.config.JWTAccessSecret.Encode(accessClaims)
	if err != nil {
		return nil, err
	}

	// Refresh token claims
	refreshClaims := map[string]interface{}{
		"type":      "refresh",
		"exp":       time.Now().Add(s.config.JWT_REFRESH_TTL).Unix(),
		"iat":       time.Now().Unix(),
		"user_uuid": userUUID.String(),
	}
	_, refreshTokenString, err := s.config.JWTRefreshSecret.Encode(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &struct {
		AccessToken  string
		RefreshToken string
	}{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthServer) ListNodesStream(req *proto.ListNodesRequest, stream proto.AuthService_ListNodesStreamServer) error {
	ctx := stream.Context()

	nodes := []*proto.Node{
		{Id: "1", Name: "Node-1", Status: "running"},
		{Id: "2", Name: "Node-2", Status: "stopped"},
		{Id: "3", Name: "Node-3", Status: "pending"},
		{Id: "4", Name: "Node-4", Status: "running"},
	}

	for _, node := range nodes {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := stream.Send(node); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
