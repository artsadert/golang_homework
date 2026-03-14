package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type DeleteUserCommand struct {
	Id uuid.UUID
}

func (u *DeleteUserCommand) Validate() error {
	if u.Id == uuid.Nil {
		return fmt.Errorf("id in User must not be empty")
	} else if uuid.Validate(u.Id.String()) != nil {
		return fmt.Errorf("uuid field must be valid")
	}
	return nil
}

type DeleteUserCommandResult struct {
	Result *common.UserResult
}
