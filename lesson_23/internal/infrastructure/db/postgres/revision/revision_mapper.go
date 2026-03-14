package revision

import "github.com/artsadert/lesson_23/internal/domain/entities"

func toDBRevision(revision *entities.Revision) *DBRevision {
	return &DBRevision{
		Uuid:       revision.Id,
		Status:     revision.Status,
		Rating:     revision.Rating,
		Review:     revision.Review,
		UserId:     revision.UserId,
		MovieId:    revision.MovieId,
		Date_added: revision.Date_added,
		Update_at:  revision.Update_at,
	}
}

func fromDBRevision(revision *DBRevision) *entities.Revision {
	return &entities.Revision{
		Id:         revision.Uuid,
		Status:     revision.Status,
		Rating:     revision.Rating,
		Review:     revision.Review,
		UserId:     revision.UserId,
		MovieId:    revision.MovieId,
		Date_added: revision.Date_added,
		Update_at:  revision.Update_at,
	}
}
