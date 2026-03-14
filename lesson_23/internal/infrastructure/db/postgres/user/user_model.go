package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBUser struct {
	gorm.Model
	Uuid      uuid.UUID `db:"uuid" gorm:"primary_key"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	Create_at time.Time `db:"create_at"`
	Update_at time.Time `db:"update_at"`
}
