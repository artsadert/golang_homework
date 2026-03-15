package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/artsadert/lesson_23/internal/application/services"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/dotenv"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/config"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/movie"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/user"
	movie_interface "github.com/artsadert/lesson_23/internal/interface/api/rest/v1/movie"
	user_interface "github.com/artsadert/lesson_23/internal/interface/api/rest/v1/user"
	"github.com/go-chi/chi"
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

	go func() {
		fmt.Println("starting on http://localhost:8080")

		err = http.ListenAndServe(":8080", root_router)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		movie_repo.Start()
	}()

	<-stop
	movie_repo.Shutdown()
	fmt.Println("shutting down")
}
