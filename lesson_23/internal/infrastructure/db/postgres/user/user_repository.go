package user

import (
	"log"

	"github.com/artsadert/lesson_23/internal/application/command"
	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repository.UserRepo {
	db.Migrator().AutoMigrate(&DBUser{})
	return &PostgresUserRepository{db: db}
}

func (u *PostgresUserRepository) Authenticate(user *command.LoginUserCommand) (*entities.User, error) {
	var db_user DBUser

	err := u.db.First(&db_user, "name = ?", user.Name).Error
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(db_user.Password), []byte(user.Password)); err != nil {
		log.Println(err)
	}

	return fromDBUser(&db_user), nil
}

func (u *PostgresUserRepository) GetUser(id uuid.UUID) (*entities.User, error) {
	var db_user DBUser

	err := u.db.First(&db_user, "uuid = ?", id).Error
	if err != nil {
		return nil, err
	}
	return fromDBUser(&db_user), nil
}

func (u *PostgresUserRepository) GetUsers() ([]*entities.User, error) {
	var db_users []*DBUser

	err := u.db.Find(&db_users).Error
	if err != nil {
		return nil, err
	}

	var users []*entities.User
	for _, db_user := range db_users {
		users = append(users, fromDBUser(db_user))
	}

	return users, nil
}

func (u *PostgresUserRepository) CreateUser(user *entities.User) error {
	db_user := toDBUser(user)

	err := u.db.Create(&db_user).Error

	return err
}

func (u *PostgresUserRepository) UpdateUser(user *entities.User) error {
	db_user := toDBUser(user)

	err := u.db.Updates(db_user).Error
	return err
}

func (u *PostgresUserRepository) DeleteUser(id uuid.UUID) error {
	var db_user DBUser
	err := u.db.Where("uuid = ?", id).Delete(&db_user).Error

	return err
}
