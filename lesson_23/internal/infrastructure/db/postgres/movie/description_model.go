package movie

import (
	"time"

	"github.com/google/uuid"
)

type DBDescriptionAggregate struct {
	Uuid       uuid.UUID `db:"uuid" gorm:"primary_key"`
	MovieId    uuid.UUID `db:"movie_id" gorm:"references:uuid;not null"`
	WordNumber int       `db:"word_number"`
	Update_at  time.Time `db:"update_at"`
	Create_at  time.Time `db:"create_at"`
}
