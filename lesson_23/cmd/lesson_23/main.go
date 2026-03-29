package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artsadert/lesson_23/internal/application/services"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/dotenv"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/config"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/movie"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/user"
	movie_interface "github.com/artsadert/lesson_23/internal/interface/api/rest/v1/movie"
	user_interface "github.com/artsadert/lesson_23/internal/interface/api/rest/v1/user"
	grpc_proto "github.com/artsadert/lesson_23/internal/interface/grpc/proto"
	grpc_middleware "github.com/artsadert/lesson_23/internal/interface/grpc/v1/interceptors"
	grpc_auth "github.com/artsadert/lesson_23/internal/interface/grpc/v1/server"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

func main() {
	dotenv.LoadDotenv()

	conn := postgres.NewConnection()
	config_repo := config.NewConfigRepo()
	config, err := config_repo.GetConfig()
	if err != nil {
		panic(err)
	}

	user_repo := user.NewPostgresUserRepository(conn)
	movie_repo := movie.NewPostgresMovieRepository(conn, 2)

	user_service := services.NewUserService(user_repo)
	movie_service := services.NewMovieService(movie_repo)

	user_mux := user_interface.NewUserMux(user_service, config)
	movie_mux := movie_interface.NewMovieMux(movie_service, config)

	root_router := chi.NewRouter()

	root_router.Mount("/users", user_mux)
	root_router.Mount("/movies", movie_mux)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":8080",
		Handler: root_router,
	}
	go func() {
		s := grpc.NewServer(
			grpc.ConnectionTimeout(5*time.Second),
			grpc.UnaryInterceptor(grpc_middleware.UnaryAuthInterceptor(config)),
			grpc.StreamInterceptor(grpc_middleware.StreamAuthInterceptor(config)),
		)

		// Register AuthService
		authServer := grpc_auth.NewAuthServer(config, user_service)
		grpc_proto.RegisterAuthServiceServer(s, authServer)

		// Register other services here...

		// Start listening
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("gRPC server listening on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		fmt.Println("starting on http://localhost:8080")

		err = server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		movie_repo.Start()
	}()

	go func() {
		movie_repo.StartErrors()
	}()

	<-stop
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	movie_repo.Shutdown(ctx)
	// server.Shutdown(ctx)
	fmt.Println("shutting down")
}
