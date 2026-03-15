package command

import (
	"github.com/artsadert/lesson_23/internal/application/common"
)

type RefreshUserCommand struct{}

func (u *RefreshUserCommand) Validate() error {
	return nil
}

type RefreshUserCommandResult struct {
	Result *common.UserResult
}
