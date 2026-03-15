package movie

import (
	"strings"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/google/uuid"
)

func toDBMovie(movie *entities.Movie) *DBMovie {
	return &DBMovie{
		Uuid:        movie.Id,
		Name:        movie.Name,
		Year:        movie.Year,
		Genre:       movie.Genre,
		Description: movie.Description,
		Poster_url:  movie.Poster_url,
		Update_at:   movie.Update_at,
		Create_at:   movie.Create_at,
	}
}

func fromDBMovie(movie *DBMovie) *entities.Movie {
	return &entities.Movie{
		Id:          movie.Uuid,
		Name:        movie.Name,
		Year:        movie.Year,
		Genre:       movie.Genre,
		Description: movie.Description,
		Poster_url:  movie.Poster_url,
		Update_at:   movie.Update_at,
		Create_at:   movie.Create_at,
	}
}

func toDBDescriptionAggregate(movie *entities.Movie) *DBDescriptionAggregate {
	return &DBDescriptionAggregate{
		Uuid:       uuid.New(),
		MovieId:    movie.Id,
		WordNumber: len(strings.Fields(movie.Description)),
	}
}
