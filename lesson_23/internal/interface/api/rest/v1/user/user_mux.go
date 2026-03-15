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
		r.Post("/login", userHandler.login)
		r.Post("/register", userHandler.createUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.DualVerifier(config))

		r.Get("/user", userHandler.getUser)
		r.Put("/users", userHandler.updateUser)
		r.Patch("/users", userHandler.updateUser)
		r.Delete("/users", userHandler.deleteUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RefreshVerifier(config))
		// r.Post("/refhresh", userHandler.regresh)
	})

	r.Group(func(r chi.Router) {
		r.Get("/users", userHandler.getUsers)
	})

	return r
}
