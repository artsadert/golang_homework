package interfaces

import (
	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/query"
	"github.com/google/uuid"
)

type RevisionService interface {
	CreateRevision(*command.CreateRevisionCommand) (*command.CreateRevisionCommandResult, error)
	DeleteRevision(*command.DeleteRevisionCommand) (*command.DeleteRevisionCommandResult, error)
	UpdateRevision(*command.UpdateRevisionCommand) (*command.UpdateRevisionCommandResult, error)
	GetRevision(uuid.UUID) (*query.RevisionQueryResult, error)
	GetRevisions() (*query.RevisionQueryListResult, error)
	GetRevisionsByUser(uuid.UUID) (*query.RevisionQueryListResult, error)
	GetRevisionsByMovie(uuid.UUID) (*query.RevisionQueryListResult, error)
}
