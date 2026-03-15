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
	db.Migrator().AutoMigrate(&DBMovie{})
	return &PostgresMovieRepository{db: db}
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

	err := p.db.Create(&dbMovie).Error
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
