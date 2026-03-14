package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/google/uuid"
)

type UserHandler struct {
	service interfaces.UserService
}

func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (userHandler *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page_size := q.Get("page_size")

	if page_size == "" {
		page_size = "10"
	}

	users, err := userHandler.service.GetUsers()
	if len(users.Result) == 0 {
		log.Println("No users found")
	}

	if err != nil {
		log.Printf("Error getting users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (userHandler *UserHandler) getUserById(w http.ResponseWriter, r *http.Request) {
	// change later to use in jwt token
	query_uuid := r.PathValue("id")
	if query_uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	err := uuid.Validate(query_uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	uuid := uuid.MustParse(query_uuid)

	// user, err := userHandler.UsersRepo.GetUserById(ctx, id)
	user, err := userHandler.service.GetUser(uuid)
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

	userCommand := &command.CreateUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userCommand.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// user, err := userHandler.UsersRepo.CreateUser(ctx, NewUser(userCommand.Name))
	userCommandResult, err := userHandler.service.CreateUser(userCommand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userCommandResult)
}

func (userHandler *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	// change to use jwt claims
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	userCommand := &command.UpdateUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userCommand.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	// err = userHandler.UsersRepo.UpdateUser(ctx, id, NewUser(userCommand.Name))
	userCommandResult, err := userHandler.service.UpdateUser(userCommand)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userCommandResult)
}

func (userHandler *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// cange later to get uuid from claims jwt tokens
	// query_uuid := r.PathValue("id")
	// if query_uuid == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }
	//
	// err := uuid.Validate(query_uuid)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// }
	//
	// uuid := uuid.MustParse(query_uuid)

	userCommand := &command.DeleteUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userCommand.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	userCommandResult, err := userHandler.service.DeleteUser(userCommand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userCommandResult)
}
