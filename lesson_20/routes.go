package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type UserHandler struct {
	Users map[int64]*User
	id    int64
	mutex sync.RWMutex
}

func NewUserHandler() *UserHandler {
	return &UserHandler{Users: map[int64]*User{}, id: 0, mutex: sync.RWMutex{}}
}

func (userHandler *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	userHandler.mutex.RLock()
	defer userHandler.mutex.RUnlock()

	json.NewEncoder(w).Encode(userHandler.Users)
}

func (userHandler *UserHandler) getUserById(w http.ResponseWriter, r *http.Request) {
	userHandler.mutex.RLock()
	defer userHandler.mutex.RUnlock()

	query_id := r.PathValue("id")
	if query_id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(query_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if _, ok := userHandler.Users[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(userHandler.Users[id])
}

func (userHandler *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	userHandler.mutex.Lock()
	defer userHandler.mutex.Unlock()

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	userCommand := CreateUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	userHandler.Users[userHandler.id] = NewUser(userCommand.Name)
	userHandler.id++

	json.NewEncoder(w).Encode(userHandler.Users[userHandler.id-1])
}

func (userHandler *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userHandler.mutex.RLock()
	defer userHandler.mutex.RUnlock()

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

	if _, ok := userHandler.Users[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userCommand := UpdateUserComamnd{}

	err = json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userCommand.Version == 0 {
		w.WriteHeader(http.StatusPreconditionRequired)
		return
	}

	if userCommand.Version != userHandler.Users[id].Version {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	userHandler.Users[id].Name = userCommand.Name
	userHandler.Users[id].Version++

	json.NewEncoder(w).Encode(userHandler.Users[id])
}

func (userHandler *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userHandler.mutex.Lock()
	defer userHandler.mutex.Unlock()

	query_id := r.PathValue("id")
	if query_id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(query_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if _, ok := userHandler.Users[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userCommand := DeleteUserCommand{}

	err = json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userCommand.Version == 0 {
		w.WriteHeader(http.StatusPreconditionRequired)
		return
	}

	if userCommand.Version != userHandler.Users[id].Version {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	delete(userHandler.Users, id)
}
