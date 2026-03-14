package common

import (
	"time"

	"github.com/google/uuid"
)

type MovieResult struct {
	Id          uuid.UUID
	Name        string
	Year        int
	Genre       []string
	Description string
	Poster_url  string
	Update_at   time.Time
	Create_at   time.Time
}
