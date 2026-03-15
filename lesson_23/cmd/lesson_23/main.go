package main

import (
	"fmt"
	"net/http"

	"github.com/artsadert/lesson_23/internal/application/services"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/dotenv"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/config"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/postgres/user"
	user_interface "github.com/artsadert/lesson_23/internal/interface/api/rest/v1/user"
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

	user_service := services.NewUserService(user_repo)

	user_mux := user_interface.NewUserMux(user_service, config)

	fmt.Println("starting on http://localhost:8080")

	err = http.ListenAndServe(":8080", user_mux)
	if err != nil {
		panic(err)
	}
}
