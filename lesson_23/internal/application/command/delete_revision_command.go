package command

import (
	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type DeleteRevisionCommand struct {
	Id uuid.UUID
}

type DeleteRevisionCommandResult struct {
	Result *common.RevisionResult
}
