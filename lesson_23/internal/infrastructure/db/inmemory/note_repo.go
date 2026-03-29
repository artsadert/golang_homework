package inmemory

import (
	"errors"
	"sync"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/google/uuid"
)

// NoteRepository is an in‑memory implementation of repository.NoteRepo.
type NoteRepository struct {
	mu    sync.RWMutex
	notes map[uuid.UUID]*entities.Note
}

// NewNoteRepository creates a new empty in‑memory note repository.
func NewNoteRepository() *NoteRepository {
	return &NoteRepository{
		notes: make(map[uuid.UUID]*entities.Note),
	}
}

// CreateNote stores a new note. The note must have a valid UUID.
func (r *NoteRepository) CreateNote(note *entities.Note) error {
	if note == nil {
		return errors.New("note cannot be nil")
	}
	if note.ID == uuid.Nil {
		note.ID = uuid.New()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.notes[note.ID]; exists {
		return errors.New("note with this ID already exists")
	}
	// Store a copy to avoid external modifications
	copyNote := *note
	r.notes[note.ID] = &copyNote
	return nil
}

// GetNote retrieves a note by its UUID.
func (r *NoteRepository) GetNote(id uuid.UUID) (*entities.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	note, exists := r.notes[id]
	if !exists {
		return nil, errors.New("note not found")
	}
	// Return a copy to prevent external modification of the stored data
	copyNote := *note
	return &copyNote, nil
}

// GetAllNotes returns all notes currently stored.
func (r *NoteRepository) GetAllNotes() ([]*entities.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	notes := make([]*entities.Note, 0, len(r.notes))
	for _, n := range r.notes {
		copyNote := *n
		notes = append(notes, &copyNote)
	}
	return notes, nil
}

// UpdateNote replaces an existing note.
func (r *NoteRepository) UpdateNote(note *entities.Note) error {
	if note == nil {
		return errors.New("note cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.notes[note.ID]; !exists {
		return errors.New("note not found")
	}
	copyNote := *note
	r.notes[note.ID] = &copyNote
	return nil
}

// DeleteNote removes a note by its UUID.
func (r *NoteRepository) DeleteNote(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.notes[id]; !exists {
		return errors.New("note not found")
	}
	delete(r.notes, id)
	return nil
}
