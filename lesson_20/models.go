package main

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewUser(name string) *User {
	return &User{Name: name}
}
