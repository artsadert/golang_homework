package repository

import (
	"github.com/artsadert/lesson_23/internal/domain/entities"
)

type ConfigRepo interface {
	GetConfig() (*entities.Config, error)
}
