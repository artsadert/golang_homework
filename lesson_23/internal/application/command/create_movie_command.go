package command

import "github.com/artsadert/lesson_23/internal/application/common"

type CreateMovieCommand struct {
	Name        string
	Year        int
	Genre       []string
	Description string
	Poster_url  string
}

type CreateMovieCommandResult struct {
	Result *common.MovieResult
}
