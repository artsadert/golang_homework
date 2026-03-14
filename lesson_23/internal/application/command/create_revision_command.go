package command

import (
	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type CreateRevisionCommand struct {
	Rating  int
	Review  string
	UserId  uuid.UUID
	MovieId uuid.UUID
}

type CreateRevisionCommandResult struct {
	Result *common.RevisionResult
}
