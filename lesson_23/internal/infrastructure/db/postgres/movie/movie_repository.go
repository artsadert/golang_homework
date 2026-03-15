package movie

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresMovieRepository struct {
	db           *gorm.DB
	workerCount  int
	taskQueue    chan uuid.UUID
	wg           *sync.WaitGroup
	shutdownFlag bool
	cancel       context.CancelFunc
}

func NewPostgresMovieRepository(db *gorm.DB, worker_count int) repository.MovieRepo {
	db.Migrator().AutoMigrate(&DBMovie{}, &DBDescriptionAggregate{})
	return &PostgresMovieRepository{
		db:           db,
		workerCount:  worker_count,
		taskQueue:    make(chan uuid.UUID, worker_count),
		wg:           &sync.WaitGroup{},
		shutdownFlag: false,
	}
}

func (p *PostgresMovieRepository) Start() error {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task_uuid := range p.taskQueue {
				movie, err := p.GetMovie(task_uuid)
				if err != nil {
					log.Printf("Error getting movie worker: %v", err)
					continue
				}
				// Чтоб проверить что рабоатает 503 ошибка когда воркеры устали от работы
				time.Sleep(20 * time.Second)

				db_description_aggregate := toDBDescriptionAggregate(movie)

				err = p.db.Create(db_description_aggregate).Error
				if err != nil {
					log.Printf("Error setting aggregate word cound worker: %v", err)
					continue
				}
			}
		}()
	}

	return nil
}

func (p *PostgresMovieRepository) upload_task(id uuid.UUID) error {
	if p.shutdownFlag {
		return fmt.Errorf("cannot upload task, repository is shutting down")
	}

	select {
	case p.taskQueue <- id:
	default:
		return fmt.Errorf("cannot upload task, queue is full")
	}

	return nil
}

func (p *PostgresMovieRepository) Shutdown(ctx context.Context) error {
	p.shutdownFlag = true

	if p.cancel != nil {
		p.cancel()
	}

	close(p.taskQueue)

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("shutdown completed without errors")
		return nil
	case <-ctx.Done():
		log.Println("shutdown completed by terminating goroutines")
		return ctx.Err()
	}
}

func (p *PostgresMovieRepository) GetMovie(id uuid.UUID) (*entities.Movie, error) {
	var dbMovie DBMovie

	err := p.db.First(&dbMovie, "uuid = ?", id).Error
	if err != nil {
		return nil, err
	}
	return fromDBMovie(&dbMovie), nil
}

func (p *PostgresMovieRepository) GetMovies() ([]*entities.Movie, error) {
	var dbMovies []*DBMovie

	err := p.db.Find(&dbMovies).Error
	if err != nil {
		return nil, err
	}

	var movies []*entities.Movie
	for _, dbMovie := range dbMovies {
		movies = append(movies, fromDBMovie(dbMovie))
	}
	return movies, nil
}

func (p *PostgresMovieRepository) CreateMovie(movie *entities.Movie) error {
	dbMovie := toDBMovie(movie)

	err := p.db.Transaction(func(tx *gorm.DB) error {
		err := p.db.Create(&dbMovie).Error
		if err != nil {
			return err
		}

		err = p.upload_task(dbMovie.Uuid)

		return err
	})

	return err
}

func (p *PostgresMovieRepository) UpdateMovie(movie *entities.Movie) error {
	dbMovie := toDBMovie(movie)

	err := p.db.Updates(dbMovie).Error
	return err
}

func (p *PostgresMovieRepository) DeleteMovie(id uuid.UUID) error {
	var dbMovie DBMovie

	err := p.db.Where("uuid = ?", id).Delete(&dbMovie).Error
	return err
}
