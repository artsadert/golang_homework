package main

import "net/http"

func CreateUserMux(userHandler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", userHandler.getUsers)

	mux.HandleFunc("/users/{id}", userHandler.getUserById)

	mux.HandleFunc("POST /users", userHandler.createUser)

	mux.HandleFunc("PUT /users/{id}", userHandler.updateUser)

	mux.HandleFunc("PATCH /users/{id}", userHandler.updateUser)

	mux.HandleFunc("DELETE /users/{id}", userHandler.deleteUser)

	return mux
}
