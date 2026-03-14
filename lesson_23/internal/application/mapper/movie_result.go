package mapper

import (
	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/artsadert/lesson_23/internal/domain/entities"
)

func NewMovieResultFromEntity(entity *entities.Movie) *common.MovieResult {
	return &common.MovieResult{Id: entity.Id,
		Name:        entity.Name,
		Year:        entity.Year,
		Genre:       entity.Genre,
		Description: entity.Description,
		Poster_url:  entity.Poster_url,
		Update_at:   entity.Update_at,
		Create_at:   entity.Create_at,
	}
}

func NewMoviesResultFromEntities(movie_entities []*entities.Movie) []*common.MovieResult {
	res := []*common.MovieResult{}

	for _, entity := range movie_entities {
		res = append(res, NewMovieResultFromEntity(entity))
	}

	return res
}
