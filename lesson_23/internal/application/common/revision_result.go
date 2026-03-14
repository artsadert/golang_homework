package common

import (
	"time"

	"github.com/google/uuid"
)

type RevisionResult struct {
	Id         uuid.UUID `json:"id"`
	Status     string    `json:"status"`
	Rating     int       `json:"raiting"`
	Review     string    `json:"review"`
	UserId     uuid.UUID `json:"user_id"`
	MovieId    uuid.UUID `json:"user_id"`
	Date_added time.Time `json:"create_at"`
	Update_at  time.Time `json:"update_at"`
}
