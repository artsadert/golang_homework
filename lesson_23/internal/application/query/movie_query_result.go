package query

import "github.com/artsadert/lesson_23/internal/application/common"

type MovieQueryResult struct {
	Result *common.MovieResult
}

type MovieQueryListResult struct {
	Result []*common.MovieResult
}
