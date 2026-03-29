package entities

import "github.com/google/uuid"

type Note struct {
	ID      uuid.UUID
	Title   string
	Content string
}
