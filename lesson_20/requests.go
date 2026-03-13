package main

type CreateUserCommand struct {
	Name string `json:"name"`
}

type UpdateUserComamnd struct {
	Name string `json:"name"`
}

type DeleteUserCommand struct{}
