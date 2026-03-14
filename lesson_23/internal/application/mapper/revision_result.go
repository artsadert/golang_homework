package mapper

import (
	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/artsadert/lesson_23/internal/domain/entities"
)

func NewRevisionResultFromEntity(entity *entities.Revision) *common.RevisionResult {
	return &common.RevisionResult{
		Id:         entity.Id,
		Status:     entity.Status,
		Rating:     entity.Rating,
		Review:     entity.Review,
		UserId:     entity.UserId,
		MovieId:    entity.MovieId,
		Date_added: entity.Date_added,
		Update_at:  entity.Update_at,
	}
}

func NewRevisionsResultFromEntities(user_entities []*entities.Revision) []*common.RevisionResult {
	res := []*common.RevisionResult{}

	for _, entity := range user_entities {
		res = append(res, NewRevisionResultFromEntity(entity))
	}

	return res
}
