package movie

import (
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/interface/api/rest/v1/user/middleware"
	"github.com/go-chi/chi/v5"
	system_middleware "github.com/go-chi/chi/v5/middleware"
)

func NewMovieMux(service interfaces.MovieService, config *entities.Config) *chi.Mux {
	movieHandler := NewMovieHandler(service)
	r := chi.NewRouter()
	r.Use(system_middleware.Logger)
	r.Use(system_middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(middleware.DualVerifier(config))

		r.Get("/movies", movieHandler.getMovies)
		r.Post("/movies", movieHandler.createMovie)
		r.Get("/movies/{id}", movieHandler.getMovie)
		r.Put("/movies/{id}", movieHandler.updateMovie)
		r.Patch("/movies/{id}", movieHandler.updateMovie)
		r.Delete("/movies/{id}", movieHandler.deleteMovie)
	})

	return r
}
