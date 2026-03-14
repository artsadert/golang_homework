package user

import (
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func NewUserMux(service interfaces.UserService, tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	userHandler := NewUserHandler(service, tokenAuth)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Post("/login", userHandler.login)
		r.Post("/register", userHandler.createUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/users", userHandler.createUser)
		r.Put("/users", userHandler.updateUser)
		r.Patch("/users", userHandler.updateUser)
		r.Delete("/users", userHandler.deleteUser)
	})

	r.Group(func(r chi.Router) {
		r.Get("/users", userHandler.getUsers)
		r.Get("/user", userHandler.getUserById)
	})

	return r
}
