package movie

import (
	"time"

	"github.com/google/uuid"
)

type DBMovie struct {
	Uuid        uuid.UUID `db:"uuid" gorm:"primary_key"`
	Name        string    `db:"name"`
	Year        int       `db:"year"`
	Genre       string    `db:"genre"`
	Description string    `db:"description"`
	Poster_url  string    `db:"poster_url"`
	Update_at   time.Time `db:"update_at"`
	Create_at   time.Time `db:"create_at"`
}
