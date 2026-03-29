package inmemory_test

import (
	"sync"
	"testing"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/infrastructure/db/inmemory"
	"github.com/google/uuid"
)

func TestConcurrentCreateNote(t *testing.T) {
	repo := inmemory.NewNoteRepository()
	const goroutines = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)

	ids := make(chan uuid.UUID, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			note := &entities.Note{
				Title:   "test",
				Content: "content",
			}
			err := repo.CreateNote(note)
			if err != nil {
				t.Errorf("CreateNote error: %v", err)
			}
			ids <- note.ID
		}()
	}

	wg.Wait()
	close(ids)

	idSet := make(map[uuid.UUID]bool)
	count := 0
	for id := range ids {
		if idSet[id] {
			t.Errorf("Duplicate ID: %v", id)
		}
		idSet[id] = true
		count++
	}
	if count != goroutines {
		t.Errorf("Expected %d notes, got %d", goroutines, count)
	}

	notes, _ := repo.GetAllNotes()
	if len(notes) != goroutines {
		t.Errorf("Repository contains %d notes, want %d", len(notes), goroutines)
	}
}
