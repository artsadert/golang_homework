package user

import (
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/interface/api/rest/v1/user/middleware"
	"github.com/go-chi/chi/v5"
	system_middleware "github.com/go-chi/chi/v5/middleware"
)

func NewUserMux(service interfaces.UserService, config *entities.Config) *chi.Mux {
	userHandler := NewUserHandler(service, config)
	r := chi.NewRouter()
	r.Use(system_middleware.Logger)
	r.Use(system_middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Post("/users/login", userHandler.login)
		r.Post("/users/register", userHandler.createUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.DualVerifier(config))

		r.Get("/users", userHandler.getUser)
		r.Put("/users", userHandler.updateUser)
		r.Patch("/users", userHandler.updateUser)
		r.Delete("/users", userHandler.deleteUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RefreshVerifier(config))
		r.Post("/users/refresh", userHandler.refresh)
	})

	return r
}
