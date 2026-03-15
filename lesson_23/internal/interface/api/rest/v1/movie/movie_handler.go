package movie

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MovieHandler struct {
	service interfaces.MovieService
}

func NewMovieHandler(service interfaces.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

// getMovies handles GET /movies
func (h *MovieHandler) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.GetMovies()
	if err != nil {
		log.Printf("Error getting movies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("Error encoding movies: %v", err)
	}
}

// getMovie handles GET /movies/{id}
func (h *MovieHandler) getMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := h.service.GetMovie(id)
	if err != nil {
		log.Printf("Error getting movie: %v", err)
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("Error encoding movie: %v", err)
	}
}

// createMovie handles POST /movies
func (h *MovieHandler) createMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var cmd command.CreateMovieCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := cmd.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	movie, err := h.service.CreateMovie(&cmd)
	if err != nil {
		if err.Error() == "cannot upload task, queue is full" {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		log.Printf("Error creating movie: %v", err)
		http.Error(w, "Internal server error", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("Error encoding movie: %v", err)
	}
}

// updateMovie handles PUT /movies/{id}
func (h *MovieHandler) updateMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var cmd command.UpdateMovieCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the ID from the URL to ensure the command targets the correct movie
	cmd.Id = id

	if err := cmd.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	movie, err := h.service.UpdateMovie(&cmd)
	if err != nil {
		log.Printf("Error updating movie: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("Error encoding movie: %v", err)
	}
}

// deleteMovie handles DELETE /movies/{id}
func (h *MovieHandler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	cmd := command.DeleteMovieCommand{Id: id}
	if err := cmd.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	movie, err := h.service.GetMovie(id)
	if err != nil {
		log.Printf("Error getting movie: %v", err)
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(movie)
}
