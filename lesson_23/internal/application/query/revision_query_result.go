package query

import "github.com/artsadert/lesson_23/internal/application/common"

type RevisionQueryResult struct {
	Result *common.RevisionResult
}

type RevisionQueryListResult struct {
	Result []*common.RevisionResult
}
