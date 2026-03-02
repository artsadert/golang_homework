package main

type CreateUserCommand struct {
	Name string `json:"name"`
}

type UpdateUserComamnd struct {
	Name    string `json:"name"`
	Version int64  `json:"version,omitempty"`
}

type DeleteUserCommand struct {
	Version int64 `json:"version,omitempty"`
}
