package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	Id          uuid.UUID
	Name        string
	Year        int
	Genre       []string
	Description string
	Poster_url  string
	Update_at   time.Time
	Create_at   time.Time
}

func NewMovie(name, description, poster_url string, year int, genre []string) (*Movie, error) {
	movie := Movie{
		Id:          uuid.New(),
		Name:        name,
		Year:        year,
		Genre:       genre,
		Description: description,
		Poster_url:  poster_url,
		Update_at:   time.Now(),
		Create_at:   time.Now(),
	}
	err := movie.validate()

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (m Movie) validate() error {
	if m.Name == "" {
		return fmt.Errorf("Name in Movie must not be empty")
	} else if m.Year <= 1890 { // first film was filmed in 1895
		return fmt.Errorf("Year must be real in Movie")
	} else if len(m.Genre) == 0 {
		return fmt.Errorf("Genre must not be empty")
	}

	for _, genre := range m.Genre {
		if genre == "" {
			return fmt.Errorf("Genre must not be empty")
		}
	}
	return nil
}

func (m *Movie) UpdateName(name string) error {
	m.Name = name
	m.Update_at = time.Now()

	return m.validate()
}

func (m *Movie) UpdateYear(year int) error {
	m.Year = year
	m.Update_at = time.Now()

	return m.validate()
}

func (m *Movie) UpdateGenre(genre []string) error {
	m.Genre = genre
	m.Update_at = time.Now()

	return m.validate()
}
