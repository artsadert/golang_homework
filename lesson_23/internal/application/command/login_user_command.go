package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
)

type LoginUserCommand struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *LoginUserCommand) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name in User must not be empty")
	} else if u.Password == "" {
		return fmt.Errorf("password in User must not be empty")
	}
	return nil
}

type LoginUserCommandResult struct {
	Result *common.UserResult
}
