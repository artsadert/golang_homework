package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
)

type CreateUserCommand struct {
	// we dont need id cause it will create it later
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *CreateUserCommand) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name in User must not be empty")
	} else if c.Email == "" {
		return fmt.Errorf("email in User must not be empty")
	}
	return nil
}

type CreateUserCommandResult struct {
	Result *common.UserResult
}
