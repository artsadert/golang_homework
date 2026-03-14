package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Create_at time.Time
	Update_at time.Time
}

func NewUser(name, email string) (*User, error) {
	user := User{
		Id:        uuid.New(),
		Name:      name,
		Email:     email,
		Update_at: time.Now(),
		Create_at: time.Now(),
	}

	err := user.validate()

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u User) validate() error {
	if u.Name == "" {
		return fmt.Errorf("Name in User must not be empty")
	} else if u.Email == "" {
		return fmt.Errorf("Email in User must not be empty")
	}

	return nil
}

func (u *User) UpdateName(name string) error {
	u.Name = name
	u.Update_at = time.Now()

	return u.validate()
}

func (u *User) UpdateEmail(email string) error {
	u.Email = email
	u.Update_at = time.Now()

	return u.validate()
}
