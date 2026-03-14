package interfaces

import (
	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/query"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(*command.CreateUserCommand) (*command.CreateUserCommandResult, error)
	DeleteUser(*command.DeleteUserCommand) (*command.DeleteUserCommandResult, error)
	UpdateUser(*command.UpdateUserCommand) (*command.UpdateUserCommandResult, error)
	GetUser(uuid.UUID) (*query.UserQueryResult, error)
	GetUsers() (*query.UserQueryListResult, error)
}
