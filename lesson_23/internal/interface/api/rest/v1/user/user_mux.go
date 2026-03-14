package user

import (
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewUserMux(service interfaces.UserService) *chi.Mux {
	userHandler := NewUserHandler(service)
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/users", userHandler.getUsers)
		r.Get("/user", userHandler.getUserById)
		r.Post("/users", userHandler.createUser)
		r.Put("/users", userHandler.updateUser)
		r.Patch("/users", userHandler.updateUser)
		r.Delete("/users", userHandler.deleteUser)
	})

	return r
}
