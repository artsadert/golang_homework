package query

import "github.com/artsadert/lesson_23/internal/application/common"

type UserQueryResult struct {
	Result *common.UserResult
}

type UserQueryListResult struct {
	Result []*common.UserResult
}
