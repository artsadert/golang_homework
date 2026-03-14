package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type UpdateUserCommand struct {
	Id    uuid.UUID `db:"uuid" gorm:"primary_key"`
	Name  string
	Email string
}

func (u *UpdateUserCommand) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name in User must not be empty")
	} else if u.Email == "" {
		return fmt.Errorf("email in User must not be empty")
	} else if u.Id == uuid.Nil {
		return fmt.Errorf("id in User must not be empty")
	} else if uuid.Validate(u.Id.String()) != nil {
		return fmt.Errorf("uuid field must be valid")
	}
	return nil
}

type UpdateUserCommandResult struct {
	Result *common.UserResult
}
