package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var revision_state = map[string]bool{
	"WATCHED": true,
	"PLANNED": true,
	"DROPPED": true,
}

type Revision struct {
	Id         uuid.UUID
	Status     string
	Rating     int
	Review     string
	UserId     uuid.UUID
	MovieId    uuid.UUID
	Date_added time.Time
	Update_at  time.Time
}

func NewRevision(rating int, review string, user_id uuid.UUID, movie_id uuid.UUID) (*Revision, error) {
	revision := Revision{
		Id:         uuid.New(),
		Status:     "PLANNED",
		Rating:     rating,
		Review:     review,
		UserId:     user_id,
		MovieId:    movie_id,
		Date_added: time.Now(),
		Update_at:  time.Now(),
	}

	err := revision.validate()

	if err != nil {
		return nil, err
	}

	return &revision, nil
}

func (r Revision) validate() error {
	if !revision_state[r.Status] {
		return fmt.Errorf("Status in Revision must be Watched, Planned or Dropped")
	} else if (r.Rating < 0 || r.Rating > 10) && r.Status == "WATCHED" { // if watched and revision not legal error
		return fmt.Errorf("Rating in Revision must be 0..10")
	}

	return nil
}

func (r *Revision) UpdateStatus(status string) error {
	r.Status = status
	r.Update_at = time.Now()

	return r.validate()
}

func (r *Revision) UpdateRating(raiting int) error {
	r.Rating = raiting
	r.Update_at = time.Now()

	return r.validate()
}

func (r *Revision) UpdateReview(review string) error {
	r.Review = review

	r.Update_at = time.Now()

	return r.validate()
}
