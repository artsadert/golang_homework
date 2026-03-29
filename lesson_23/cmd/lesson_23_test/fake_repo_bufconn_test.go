package lesson_23_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/mapper"
	"github.com/artsadert/lesson_23/internal/application/query"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/artsadert/lesson_23/internal/interface/grpc/proto"
	grpc_auth "github.com/artsadert/lesson_23/internal/interface/grpc/v1/interceptors"
	"github.com/artsadert/lesson_23/internal/interface/grpc/v1/server"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// ---------- Fake in‑memory user repository ----------
type fakeUserRepo struct {
	users map[string]*entities.User // key = name
}

func newFakeUserRepo() *fakeUserRepo {
	// Pre‑register user "hello" with password "test" (bcrypt hashed)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := &entities.User{
		Id:       uuid.New(),
		Name:     "hello",
		Password: string(hashed),
	}
	return &fakeUserRepo{
		users: map[string]*entities.User{user.Name: user},
	}
}

func (r *fakeUserRepo) Authenticate(cmd *command.LoginUserCommand) (*entities.User, error) {
	user, ok := r.users[cmd.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cmd.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}
	return user, nil
}

func (r *fakeUserRepo) GetUser(id uuid.UUID) (*entities.User, error) {
	for _, u := range r.users {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, status.Error(codes.NotFound, "user not found")
}

// Other methods (GetUsers, CreateUser, etc.) – implement if needed, but not required for auth.
func (r *fakeUserRepo) GetUsers() ([]*entities.User, error)  { return nil, nil }
func (r *fakeUserRepo) CreateUser(user *entities.User) error { return nil }
func (r *fakeUserRepo) UpdateUser(user *entities.User) error { return nil }
func (r *fakeUserRepo) DeleteUser(id uuid.UUID) error        { return nil }

// ---------- Fake JWT config ----------
func testConfig() *entities.Config {
	accessSecret := jwtauth.New("HS256", []byte("test-access-secret"), nil) // implement if needed
	refreshSecret := jwtauth.New("HS256", []byte("test-refresh-secret"), nil)
	return &entities.Config{
		JWTAccessSecret:  accessSecret,
		JWTRefreshSecret: refreshSecret,
		JWT_TTL:          15 * time.Minute,
		JWT_REFRESH_TTL:  24 * time.Hour,
	}
}

// ---------- The integration test ----------
func TestAuthServiceWithBufconn(t *testing.T) {
	// 1. Create fake dependencies
	userRepo := newFakeUserRepo()
	userService := &fakeUserService{repo: userRepo} // adapt to your UserService interface
	cfg := testConfig()

	// 2. Create the real AuthServer (with fake repo & config)
	authServer := server.NewAuthServer(cfg, userService)

	// 3. Set up gRPC server with auth interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryAuthInterceptor(cfg)),
		grpc.StreamInterceptor(grpc_auth.StreamAuthInterceptor(cfg)),
	)
	proto.RegisterAuthServiceServer(grpcServer, authServer)

	// 4. Create bufconn listener and start server
	listener := bufconn.Listen(bufSize)
	defer listener.Close()
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Logf("bufconn server exited: %v", err)
		}
	}()

	// 5. Create client connection via bufconn
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewAuthServiceClient(conn)

	// ---------- 6. Test login ----------
	loginCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	loginResp, err := client.Login(loginCtx, &proto.LoginRequest{
		Name:     "hello",
		Password: "test",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	validToken := loginResp.AccessToken
	t.Log("Login successful, token obtained")

	// ---------- 7. Test streaming with valid token (short timeout) ----------
	t.Run("Valid token + timeout", func(t *testing.T) {
		ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		md := metadata.Pairs("authorization", "Bearer "+validToken)
		ctx2 = metadata.NewOutgoingContext(ctx2, md)

		stream, err := client.ListNodesStream(ctx2, &proto.ListNodesRequest{})
		if err != nil {
			t.Fatalf("Failed to start stream: %v", err)
		}

		var received int
		for {
			_, err := stream.Recv()
			if err != nil {
				if status.Code(err) == codes.DeadlineExceeded {
					t.Logf("DeadlineExceeded after %d nodes – as expected", received)
					break
				}
				t.Fatalf("Unexpected error: %v", err)
			}
			received++
		}
	})

	// ---------- 8. Test invalid token ----------
	t.Run("Invalid token", func(t *testing.T) {
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		md := metadata.Pairs("authorization", "Bearer invalid.token.here")
		ctx2 = metadata.NewOutgoingContext(ctx2, md)

		stream, err := client.ListNodesStream(ctx2, &proto.ListNodesRequest{})
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

	// ---------- 9. Test no token ----------
	t.Run("No token", func(t *testing.T) {
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		// No metadata attached
		stream, err := client.ListNodesStream(ctx2, &proto.ListNodesRequest{})
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
}

// Adapter to satisfy interfaces.UserService
type fakeUserService struct {
	repo repository.UserRepo
}

func (s *fakeUserService) Authenticate(cmd *command.LoginUserCommand) (*query.UserQueryResult, error) {
	user, err := s.repo.Authenticate(cmd)
	if err != nil {
		return nil, err
	}
	return &query.UserQueryResult{Result: mapper.NewUserResultFromEntity(user)}, nil
}

func (s *fakeUserService) GetUser(id uuid.UUID) (*query.UserQueryResult, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}
	return &query.UserQueryResult{Result: mapper.NewUserResultFromEntity(user)}, nil
}

func (s *fakeUserService) GetUsers() (*query.UserQueryListResult, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return &query.UserQueryListResult{Result: mapper.NewUsersResultFromEntities(users)}, nil
}

func (s *fakeUserService) CreateUser(*command.CreateUserCommand) (*command.CreateUserCommandResult, error) {
	return nil, nil
}

func (s *fakeUserService) UpdateUser(*command.UpdateUserCommand) (*command.UpdateUserCommandResult, error) {
	return nil, nil
}

func (s *fakeUserService) DeleteUser(*command.DeleteUserCommand) (*command.DeleteUserCommandResult, error) {
	return nil, nil
}
