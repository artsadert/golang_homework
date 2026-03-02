package main

import (
	"net/http"
)

func main() {
	userHandler := NewUserHandler()

	mux := CreateUserMux(userHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
