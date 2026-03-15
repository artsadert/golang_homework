package command

import (
	"fmt"

	"github.com/artsadert/lesson_23/internal/application/common"
)

type CreateMovieCommand struct {
	Name        string
	Year        int
	Genre       string
	Description string
	Poster_url  string
}

func (c *CreateMovieCommand) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name in Movie must not be empty")
	} else if c.Year <= 1890 { // first film was filmed in 1895
		return fmt.Errorf("year must be real in Movie")
	} else if c.Genre == "" {
		return fmt.Errorf("genre must not be empty")
	} else if c.Description == "" {
		return fmt.Errorf("description must not be empty")
	}
	return nil
}

type CreateMovieCommandResult struct {
	Result *common.MovieResult
}
