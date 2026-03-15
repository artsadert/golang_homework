package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service interfaces.UserService
	config  *entities.Config
}

func NewUserHandler(service interfaces.UserService, config *entities.Config) *UserHandler {
	return &UserHandler{
		service: service,
		config:  config,
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

	tokens, err := createTokens(user.Result.Id, userHandler.config)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(tokens)
}

func (userHandler *UserHandler) refresh(w http.ResponseWriter, r *http.Request) {
	RefreshUserCommand := &command.RefreshUserCommand{}

	err := json.NewDecoder(r.Body).Decode(RefreshUserCommand)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = RefreshUserCommand.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query_uuid, err := getUUIDFromContextToken(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := userHandler.service.GetUser(query_uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tokens, err := refreshToken(user.Result.Id, userHandler.config)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(tokens)
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

func (userHandler *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	query_uuid, err := getUUIDFromContextToken(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := userHandler.service.GetUser(query_uuid)
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

	tokens, err := createTokens(userCommandResult.Result.Id, userHandler.config)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode([]interface{}{userCommandResult, tokens})
}

func (userHandler *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	userCommand := &command.UpdateUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if claims["user_uuid"] == nil {
		http.Error(w, "user_uuid must by in token", http.StatusUnauthorized)
	}
	userCommand.Id = uuid.MustParse(claims["user_uuid"].(string))

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
	userCommand := &command.DeleteUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if claims["user_uuid"] == nil {
		http.Error(w, "user_uuid must by in token", http.StatusUnauthorized)
	}
	userCommand.Id = uuid.MustParse(claims["user_uuid"].(string))

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
