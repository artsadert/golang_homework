package command

import (
	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type DeleteMovieCommand struct {
	Id uuid.UUID
}

type DeleteMovieCommandResult struct {
	Result *common.MovieResult
}
