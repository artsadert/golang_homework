package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	UsersRepo *UserRepo
}

func NewUserHandler(userRepo *UserRepo) *UserHandler {
	return &UserHandler{UsersRepo: userRepo}
}

func (userHandler *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page_size := q.Get("page_size")

	if page_size == "" {
		page_size = "10"
	}

	page, err := strconv.Atoi(page_size)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	cursor := q.Get("cursor")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	users, err := userHandler.UsersRepo.GetUsers(ctx, page, cursor)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (userHandler *UserHandler) getUserById(w http.ResponseWriter, r *http.Request) {
	query_id := r.PathValue("id")
	if query_id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(query_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	user, err := userHandler.UsersRepo.GetUserById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (userHandler *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	userCommand := CreateUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	user, err := userHandler.UsersRepo.CreateUser(ctx, NewUser(userCommand.Name))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (userHandler *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	query_id := r.PathValue("id")
	if query_id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(query_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	userCommand := UpdateUserComamnd{}

	err = json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = userHandler.UsersRepo.UpdateUser(ctx, id, NewUser(userCommand.Name))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (userHandler *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	query_id := r.PathValue("id")
	if query_id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(query_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	userCommand := DeleteUserCommand{}

	err = json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = userHandler.UsersRepo.DeleteUser(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
