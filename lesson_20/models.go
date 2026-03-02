package main

type User struct {
	Name    string `json:"name"`
	Version int64  `json:"version"`
}

func NewUser(name string) *User {
	return &User{Name: name, Version: 1}
}
