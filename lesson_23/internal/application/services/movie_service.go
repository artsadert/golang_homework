package services

import (
	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/application/mapper"
	"github.com/artsadert/lesson_23/internal/application/query"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
)

type MovieService struct {
	movieRepo repository.MovieRepo
}

func NewMovieService(repo repository.MovieRepo) interfaces.MovieService {
	return &MovieService{movieRepo: repo}
}

func (m MovieService) GetMovie(id uuid.UUID) (*query.MovieQueryResult, error) {
	entity, err := m.movieRepo.GetMovie(id)
	if err != nil {
		return nil, err
	}

	return &query.MovieQueryResult{Result: mapper.NewMovieResultFromEntity(entity)}, nil
}

func (m MovieService) GetMovies() (*query.MovieQueryListResult, error) {
	movie_entities, err := m.movieRepo.GetMovies()
	if err != nil {
		return nil, err
	}

	return &query.MovieQueryListResult{Result: mapper.NewMoviesResultFromEntities(movie_entities)}, nil
}

func (m MovieService) CreateMovie(movie *command.CreateMovieCommand) (*command.CreateMovieCommandResult, error) {
	entity, err := entities.NewMovie(movie.Name, movie.Description, movie.Poster_url, movie.Year, movie.Genre)
	if err != nil {
		return nil, err
	}

	err = m.movieRepo.CreateMovie(entity)
	if err != nil {
		return nil, err
	}
	return &command.CreateMovieCommandResult{Result: mapper.NewMovieResultFromEntity(entity)}, nil
}

func (m MovieService) DeleteMovie(movie *command.DeleteMovieCommand) (*command.DeleteMovieCommandResult, error) {
	entity, err := m.movieRepo.GetMovie(movie.Id)
	if err != nil {
		return nil, err
	}

	err = m.movieRepo.DeleteMovie(movie.Id)
	if err != nil {
		return nil, err
	}

	return &command.DeleteMovieCommandResult{Result: mapper.NewMovieResultFromEntity(entity)}, nil
}

func (m MovieService) UpdateMovie(movie *command.UpdateMovieCommand) (*command.UpdateMovieCommandResult, error) {
	entity, err := m.movieRepo.GetMovie(movie.Id)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateName(movie.Name)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateYear(movie.Year)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateGenre(movie.Genre)
	if err != nil {
		return nil, err
	}

	return &command.UpdateMovieCommandResult{Result: mapper.NewMovieResultFromEntity(entity)}, nil
}
