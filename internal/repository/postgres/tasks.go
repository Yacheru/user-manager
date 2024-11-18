package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user-manager/internal/entities"
)

type Tasks struct {
	db *sqlx.DB
}

func NewTasks(db *sqlx.DB) *Tasks {
	return &Tasks{db: db}
}

func (t *Tasks) NewTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	query := `
		INSERT INTO tasks (uuid, title, description, reward) 
		VALUES ($1, $2, $3, $4)
	`
	if _, err := t.db.ExecContext(ctx, query, task.TaskID, task.Title, task.Description, task.Reward); err != nil {
		return nil, err
	}

	return task, nil
}
