package revision

import (
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRevisionRepo struct {
	db *gorm.DB
}

func NewPostgresRevisionRepo(db *gorm.DB) repository.RevisionRepo {
	return &PostgresRevisionRepo{db: db}
}

func (r *PostgresRevisionRepo) GetRevision(id uuid.UUID) (*entities.Revision, error) {
	var db_revision DBRevision
	err := r.db.Where("uuid = ?", id).First(&db_revision).Error
	if err != nil {
		return nil, err
	}
	return fromDBRevision(&db_revision), nil
}

func (r *PostgresRevisionRepo) GetRevisions() ([]*entities.Revision, error) {
	var db_revisions []*DBRevision

	err := r.db.Take(&db_revisions).Error
	if err != nil {
		return nil, err
	}

	var revisions []*entities.Revision
	for _, db_revision := range db_revisions {
		revisions = append(revisions, fromDBRevision(db_revision))
	}

	return revisions, nil
}

func (r *PostgresRevisionRepo) GetRevisionsByUserId(id uuid.UUID) ([]*entities.Revision, error) {
	var db_revisions []*DBRevision

	err := r.db.Where("user_id = ?", id).Take(&db_revisions).Error
	if err != nil {
		return nil, err
	}

	var revisions []*entities.Revision
	for _, db_revision := range db_revisions {
		revisions = append(revisions, fromDBRevision(db_revision))
	}

	return revisions, nil
}

func (r *PostgresRevisionRepo) GetRevisionsByMovieId(id uuid.UUID) ([]*entities.Revision, error) {
	var db_revisions []*DBRevision

	err := r.db.Where("movie_id = ?", id).Take(&db_revisions).Error
	if err != nil {
		return nil, err
	}

	var revisions []*entities.Revision
	for _, db_revision := range db_revisions {
		revisions = append(revisions, fromDBRevision(db_revision))
	}

	return revisions, nil
}

func (r *PostgresRevisionRepo) CreateRevision(revision *entities.Revision) error {
	db_revision := toDBRevision(revision)

	err := r.db.Create(db_revision).Error
	return err
}

func (r *PostgresRevisionRepo) UpdateRevision(revision *entities.Revision) error {
	db_revision := toDBRevision(revision)

	err := r.db.Save(db_revision).Error

	return err
}

func (r *PostgresRevisionRepo) DeleteRevision(id uuid.UUID) error {
	var db_movie DBRevision

	err := r.db.Where("uuid = ?", id).Delete(&db_movie).Error

	return err
}
