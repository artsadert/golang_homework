package services

import (
	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/application/mapper"
	"github.com/artsadert/lesson_23/internal/application/query"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
)

type RevisionService struct {
	revisionRepo repository.RevisionRepo
}

func NewRevisionService(repo repository.RevisionRepo) interfaces.RevisionService {
	return &RevisionService{revisionRepo: repo}
}

func (r RevisionService) GetRevision(id uuid.UUID) (*query.RevisionQueryResult, error) {
	entity, err := r.revisionRepo.GetRevision(id)
	if err != nil {
		return nil, err
	}
	return &query.RevisionQueryResult{Result: mapper.NewRevisionResultFromEntity(entity)}, nil
}

func (r RevisionService) GetRevisions() (*query.RevisionQueryListResult, error) {
	revision_entities, err := r.revisionRepo.GetRevisions()
	if err != nil {
		return nil, err
	}

	return &query.RevisionQueryListResult{Result: mapper.NewRevisionsResultFromEntities(revision_entities)}, nil
}

func (r RevisionService) GetRevisionsByUser(id uuid.UUID) (*query.RevisionQueryListResult, error) {
	revision_entities, err := r.revisionRepo.GetRevisionsByUserId(id)
	if err != nil {
		return nil, err
	}

	return &query.RevisionQueryListResult{Result: mapper.NewRevisionsResultFromEntities(revision_entities)}, nil
}

func (r RevisionService) GetRevisionsByMovie(id uuid.UUID) (*query.RevisionQueryListResult, error) {
	revision_entities, err := r.revisionRepo.GetRevisionsByMovieId(id)
	if err != nil {
		return nil, err
	}

	return &query.RevisionQueryListResult{Result: mapper.NewRevisionsResultFromEntities(revision_entities)}, nil
}

func (r RevisionService) CreateRevision(revision *command.CreateRevisionCommand) (*command.CreateRevisionCommandResult, error) {
	entity, err := entities.NewRevision(revision.Rating, revision.Review, revision.UserId, revision.MovieId)
	if err != nil {
		return nil, err
	}

	err = r.revisionRepo.CreateRevision(entity)
	if err != nil {
		return nil, err
	}

	return &command.CreateRevisionCommandResult{Result: mapper.NewRevisionResultFromEntity(entity)}, nil
}

func (r RevisionService) DeleteRevision(revision *command.DeleteRevisionCommand) (*command.DeleteRevisionCommandResult, error) {
	entity, err := r.revisionRepo.GetRevision(revision.Id)
	if err != nil {
		return nil, err
	}

	err = r.revisionRepo.DeleteRevision(revision.Id)
	if err != nil {
		return nil, err
	}

	return &command.DeleteRevisionCommandResult{Result: mapper.NewRevisionResultFromEntity(entity)}, nil
}

func (r RevisionService) UpdateRevision(revision *command.UpdateRevisionCommand) (*command.UpdateRevisionCommandResult, error) {
	entity, err := r.revisionRepo.GetRevision(revision.Id)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateRating(revision.Rating)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateReview(revision.Review)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateStatus(revision.Status)
	if err != nil {
		return nil, err
	}

	return &command.UpdateRevisionCommandResult{Result: mapper.NewRevisionResultFromEntity(entity)}, nil
}
