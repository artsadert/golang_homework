package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
	"github.com/google/uuid"
)

type UpdateMovieCommand struct {
	Id          uuid.UUID
	Name        string
	Year        int
	Genre       string
	Description string
	Poster_url  string
}

func (c *UpdateMovieCommand) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name in User must not be empty")
	} else if c.Year <= 1890 { // first film was filmed in 1895
		return fmt.Errorf("year must be real in Movie")
	} else if len(c.Genre) == 0 {
		return fmt.Errorf("genre must not be empty")
	}
	return nil
}

type UpdateMovieCommandResult struct {
	Result *common.MovieResult
}
