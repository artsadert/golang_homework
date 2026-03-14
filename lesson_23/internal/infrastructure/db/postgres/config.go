package postgres

import (
	"fmt"
	"os"
)

func NewPostgresDSN() string {
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	ssl_mode := os.Getenv("DB_SSLMODE")

	db_url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, username, password, db_name, port, ssl_mode)

	return db_url
}
