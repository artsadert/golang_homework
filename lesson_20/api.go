package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateUserMux(userHandler *UserHandler) *chi.Mux {
	r := chi.NewRouter()

	// Add some useful middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes (if any)
	r.Group(func(r chi.Router) {
		// Add auth middleware to protect user routes
		r.Use(AuthorizeUserAccess)

		r.Get("/users", userHandler.getUsers)
		r.Get("/users/{id}", userHandler.getUserById)
		r.Post("/users", userHandler.createUser)
		r.Put("/users/{id}", userHandler.updateUser)
		r.Patch("/users/{id}", userHandler.updateUser)
		r.Delete("/users/{id}", userHandler.deleteUser)
	})

	return r
}
