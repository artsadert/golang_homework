package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service   interfaces.UserService
	tokenAuth *jwtauth.JWTAuth
}

func NewUserHandler(service interfaces.UserService, tokenAuth *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{
		service:   service,
		tokenAuth: tokenAuth,
	}
}

func (userHandler *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	LoginUserCommand := &command.LoginUserCommand{}

	err := json.NewDecoder(r.Body).Decode(LoginUserCommand)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = LoginUserCommand.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate credentials against your user service
	user, err := userHandler.service.Authenticate(LoginUserCommand)
	if err != nil || user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, tokenString, err := userHandler.tokenAuth.Encode(map[string]interface{}{"user_id": user.Result.Id})
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
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
		log.Fatalf("Failed to create user: %v", err)
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

	err := json.NewDecoder(r.Body).Decode(userCommand)
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

	err := json.NewDecoder(r.Body).Decode(userCommand)
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
