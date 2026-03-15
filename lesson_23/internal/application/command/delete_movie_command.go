package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type DeleteMovieCommand struct {
	Id uuid.UUID
}

func (c *DeleteMovieCommand) Validate() error {
	if c.Id == uuid.Nil {
		return fmt.Errorf("id in User must not be empty")
	} else if uuid.Validate(c.Id.String()) != nil {
		return fmt.Errorf("uuid field must be valid")
	}
	return nil
}

type DeleteMovieCommandResult struct {
	Result *common.MovieResult
}
