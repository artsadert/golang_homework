package movie

import (
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresMovieRepository struct {
	db *gorm.DB
}

func NewPostgresMovieRepository(db *gorm.DB) repository.MovieRepo {
	return &PostgresMovieRepository{db: db}
}

func (p *PostgresMovieRepository) GetMovie(id uuid.UUID) (*entities.Movie, error) {
	var db_movie DBMovie

	err := p.db.Where("uuid = ?", id).First(&db_movie).Error
	if err != nil {
		return nil, err
	}
	return fromDBMovie(&db_movie), nil
}

func (p *PostgresMovieRepository) GetMovies() ([]*entities.Movie, error) {
	var db_movies []*DBMovie

	err := p.db.Take(&db_movies).Error
	if err != nil {
		return nil, err
	}

	var movies []*entities.Movie
	for _, db_movie := range db_movies {
		movies = append(movies, fromDBMovie(db_movie))
	}

	return movies, nil
}

func (p *PostgresMovieRepository) CreateMovie(movie *entities.Movie) error {
	db_movie := toDBMovie(movie)

	err := p.db.Create(db_movie).Error

	return err
}

func (m *PostgresMovieRepository) UpdateMovie(movie *entities.Movie) error {
	db_movie := toDBMovie(movie)

	err := m.db.Save(db_movie).Error

	return err
}

func (m *PostgresMovieRepository) DeleteMovie(id uuid.UUID) error {
	var db_movie DBMovie

	err := m.db.Where("uuid = ?", id).Delete(&db_movie).Error

	return err
}
