package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type LogRepo struct {
	db *sql.DB
}

func NewLogRepo(db *sql.DB) *LogRepo {
	return &LogRepo{db: db}
}

func (logRepo *LogRepo) CreateEntityTable(ctx context.Context, db *sql.DB, tableName string) error {
	err := logRepo.addLog(ctx, fmt.Sprintf("created %s", tableName))
	return err
}

func (logRepo *LogRepo) UpdatedEntityTable(ctx context.Context, db *sql.DB, tableName string) error {
	err := logRepo.addLog(ctx, fmt.Sprintf("updated %s", tableName))
	return err
}

func (logRepo *LogRepo) DeleteEntityTable(ctx context.Context, db *sql.DB, tableName string) error {
	err := logRepo.addLog(ctx, fmt.Sprintf("deleted %s", tableName))
	return err
}

func (logRepo *LogRepo) addLog(ctx context.Context, log_message string) error {
	query := `
			INSERT INTO logs (log_level, message, created_at)
				VALUES ($1, $2, $3)
			`

	_, err := logRepo.db.ExecContext(ctx, query, "info", log_message, time.Now())
	if err != nil {
		log.Printf("Error adding log: %v", err)
	}

	return err
}

func (logRepo *LogRepo) GetLogs() []string {
	query := `
			SELECT message
			FROM logs
			ORDER BY created_at DESC
		`

	rows, err := logRepo.db.Query(query)
	if err != nil {
		log.Printf("Error getting logs: %v", err)
	}

	var logs []string
	for rows.Next() {
		var message string
		err := rows.Scan(&message)
		if err != nil {
			log.Printf("Error scanning log: %v", err)
		}
		logs = append(logs, message)
	}
	return logs
}
