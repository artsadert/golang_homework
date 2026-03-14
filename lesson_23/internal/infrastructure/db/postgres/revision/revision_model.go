package revision

import (
	"time"

	"github.com/google/uuid"
)

type DBRevision struct {
	Uuid       uuid.UUID `db:"uuid"`
	Status     string    `db:"status"`
	Rating     int       `db:"raiting"`
	Review     string    `db:"review"`
	UserId     uuid.UUID `db:"user_id"`
	MovieId    uuid.UUID `db:"movie_id"`
	Date_added time.Time `db:"create_at"`
	Update_at  time.Time `db:"update_at"`
}
