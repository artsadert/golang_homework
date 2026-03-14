package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(NewPostgresDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
