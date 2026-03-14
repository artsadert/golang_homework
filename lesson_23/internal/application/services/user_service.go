package services

import (
	"log"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/application/interfaces"
	"github.com/artsadert/lesson_23/internal/application/mapper"
	"github.com/artsadert/lesson_23/internal/application/query"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) interfaces.UserService {
	return &UserService{userRepo: repo}
}

func (u UserService) GetUser(id uuid.UUID) (*query.UserQueryResult, error) {
	entity, err := u.userRepo.GetUser(id)
	if err != nil {
		return nil, err
	}
	return &query.UserQueryResult{Result: mapper.NewUserResultFromEntity(entity)}, nil
}

func (u UserService) GetUsers() (*query.UserQueryListResult, error) {
	user_entities, err := u.userRepo.GetUsers()
	if err != nil {
		return nil, err
	}

	return &query.UserQueryListResult{Result: mapper.NewUsersResultFromEntities(user_entities)}, nil
}

func (u UserService) CreateUser(user *command.CreateUserCommand) (*command.CreateUserCommandResult, error) {
	entity, err := entities.NewUser(user.Name, user.Email)
	if err != nil {
		return nil, err
	}

	log.Println("hello")
	err = u.userRepo.CreateUser(entity)
	if err != nil {
		return nil, err
	}

	return &command.CreateUserCommandResult{Result: mapper.NewUserResultFromEntity(entity)}, nil
}

func (u UserService) DeleteUser(user *command.DeleteUserCommand) (*command.DeleteUserCommandResult, error) {
	entity, err := u.userRepo.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	err = u.userRepo.DeleteUser(user.Id)
	if err != nil {
		return nil, err
	}
	return &command.DeleteUserCommandResult{Result: mapper.NewUserResultFromEntity(entity)}, nil
}

func (u UserService) UpdateUser(user *command.UpdateUserCommand) (*command.UpdateUserCommandResult, error) {
	entity, err := u.userRepo.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateEmail(user.Email)
	if err != nil {
		return nil, err
	}

	err = entity.UpdateName(user.Name)
	if err != nil {
		return nil, err
	}

	err = u.userRepo.UpdateUser(entity)
	if err != nil {
		return nil, err
	}

	return &command.UpdateUserCommandResult{Result: mapper.NewUserResultFromEntity(entity)}, nil
}
