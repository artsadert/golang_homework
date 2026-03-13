package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	LogRepo := NewLogRepo(db)
	UserRepo := NewUserRepo(db, LogRepo)

	userHandler := NewUserHandler(UserRepo)

	mux := CreateUserMux(userHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
