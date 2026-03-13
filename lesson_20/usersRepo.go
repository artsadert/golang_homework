package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type UserRepo struct {
	db      *sql.DB
	logRepo *LogRepo
}

type PostgresUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserRepo(db *sql.DB, logRepo *LogRepo) *UserRepo {
	return &UserRepo{db: db, logRepo: logRepo}
}

// crud operations

func (userRepo *UserRepo) GetUsers(ctx context.Context, page int, cursor string) ([]PostgresUser, error) {
	query := `
				SELECT id, name, created_at, updated_at
				FROM users
				order by id desc
				limit $1
			`

	if cursor != "" {
		splited := strings.Split(cursor, "|")
		cursor_date, err := time.Parse(time.RFC3339, splited[0])
		if err != nil {
			return nil, err
		}
		cursor_id, err := strconv.Atoi(splited[1])
		if err != nil {
			return nil, err
		}

		query = `
				SELECT id, name, created_at, updated_at
				FROM users
				where created_at > $2 and id > $3
				order by id desc
				limit $1
			`

		rows, err := userRepo.db.QueryContext(ctx, query, page, cursor_date, cursor_id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var users []PostgresUser
		for rows.Next() {
			var user PostgresUser
			err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		return users, nil
	}

	rows, err := userRepo.db.QueryContext(ctx, query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []PostgresUser
	for rows.Next() {
		var user PostgresUser
		err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (userRepo *UserRepo) GetUserById(ctx context.Context, id int64) (PostgresUser, error) {
	query := `
				SELECT id, name, created_at, updated_at
				FROM users
				WHERE id = $1
			`

	var user PostgresUser
	err := userRepo.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (userRepo *UserRepo) CreateUser(ctx context.Context, user *User) (*PostgresUser, error) {
	tx, err := userRepo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO users (name, created_at, updated_at)
        VALUES ($1, $2, $3)
        RETURNING id, name,  created_at, updated_at
    `

	var newUser PostgresUser
	now := time.Now()

	err = userRepo.db.QueryRowContext(ctx, query, user.Name, now, now).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	err = userRepo.logRepo.CreateEntityTable(ctx, userRepo.db, fmt.Sprintf("users %d", newUser.ID))
	if err != nil {
		log.Printf("Error creating entity table: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	return &newUser, nil
}

func (userRepo *UserRepo) UpdateUser(ctx context.Context, id int64, user *User) error {
	tx, err := userRepo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	query := `
				UPDATE users
				SET name = $1, updated_at = $2
				WHERE id = $3
			`

	res, err := userRepo.db.ExecContext(ctx, query, user.Name, time.Now(), id)
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	if rows == 0 {
		log.Printf("Total rows which was updated is 0")
		return nil
	}

	err = userRepo.logRepo.UpdatedEntityTable(ctx, userRepo.db, fmt.Sprintf("users %d", id))
	if err != nil {
		log.Printf("Error updating entity table: %v", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return err
}

func (userRepo *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	tx, err := userRepo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	query := `
				DELETE FROM users
				WHERE id = $1
			`

	_, err = userRepo.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	err = userRepo.logRepo.DeleteEntityTable(ctx, userRepo.db, fmt.Sprintf("users %d", id))
	if err != nil {
		log.Printf("Error deleting entity table: %v", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	return err
}
