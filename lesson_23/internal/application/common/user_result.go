package common

import (
	"time"

	"github.com/google/uuid"
)

type UserResult struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Create_at time.Time `json:"create_at"`
	Update_at time.Time `json":update_at"`
}
