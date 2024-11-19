package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"

	"user-manager/init/logger"
	"user-manager/internal/entities"
	"user-manager/pkg/constants"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}

func (u *User) NewUser(ctx context.Context, user *entities.NewUser) (*entities.User, error) {
	var userEntity = new(entities.User)

	query := `
		INSERT INTO users (uuid, name) VALUES ($1, $2) RETURNING uuid, name, points, referral
	`
	if err := u.db.GetContext(ctx, userEntity, query, user.UserID, user.Name); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, constants.UserAlreadyExistError
		}
		logger.Error(err.Error(), constants.PostgresCategory)
		return nil, err
	}

	return userEntity, nil
}

func (u *User) GetStatus(ctx context.Context, userId string) (*entities.User, error) {
	query := `
		SELECT u.uuid, u.name, u.points, u.referral,
		   	   t.uuid, t.title, t.description, t.reward
		FROM users u
			LEFT JOIN user_completed_tasks ct ON u.id = ct.user_id
			LEFT JOIN tasks t ON ct.task_id = t.id
		WHERE u.uuid = $1
	`
	rows, err := u.db.QueryxContext(ctx, query, userId)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return nil, err
	}
	defer rows.Close()

	var userEntity = new(entities.User)

	for rows.Next() {
		var task = new(entities.Task)

		if err := rows.Scan(
			&userEntity.UserID, &userEntity.Name, &userEntity.Points, &userEntity.Referral,
			&task.TaskID, &task.Title, &task.Description, &task.Reward,
		); err != nil {
			if strings.Contains(err.Error(), "converting NULL to string is unsupported") {
				logger.DebugF("no tasks found for user: %s", constants.PostgresCategory, userId)
			} else {
				logger.Error(err.Error(), constants.PostgresCategory)
				return nil, err
			}
		}
		if task.TaskID != nil {
			userEntity.Tasks = append(userEntity.Tasks, *task)
		}
	}

	if len(userEntity.Tasks) == 0 {
		userEntity.Tasks = []entities.Task{}
	}

	if userEntity.UserID != userId {
		return nil, constants.UserNotFoundError
	}

	return userEntity, nil
}

func (u *User) Leaderboard(ctx context.Context, limit, offset string) (*[]entities.Leaderboard, error) {
	var userEntity = new([]entities.Leaderboard)
	query := `
		SELECT uuid, name, points, referral
		FROM users 
		ORDER BY points DESC
		LIMIT $1
		OFFSET $2
	`
	if err := u.db.SelectContext(ctx, userEntity, query, limit, offset); err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return nil, err
	}

	return userEntity, nil
}

func (u *User) TaskComplete(ctx context.Context, taskId, userId string) error {
	var taskEntity = new(entities.Task)
	taskQuery := `
		SELECT id, uuid, title, description, reward FROM tasks WHERE uuid = $1
	`
	if err := u.db.GetContext(ctx, taskEntity, taskQuery, taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.TaskNotFoundError
		}
		logger.Error("task query: "+err.Error(), constants.PostgresCategory)
		return err
	}

	var userEntity = new(entities.User)
	userQuery := `
		SELECT id, uuid, name, points, referral FROM users WHERE uuid = $1
	`
	if err := u.db.GetContext(ctx, userEntity, userQuery, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.UserNotFoundError
		}
		logger.Error("user query: "+err.Error(), constants.PostgresCategory)
		return err
	}

	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO user_completed_tasks (user_id, task_id) 
		SELECT $1, $2
    	WHERE NOT EXISTS (
        	SELECT 1 FROM user_completed_tasks WHERE user_id = $1 AND task_id = $2
    	)
	`
	res, err := tx.ExecContext(ctx, insertQuery, userEntity.ID, taskEntity.ID)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		tx.Rollback()
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return constants.TaskAlreadyCompletedError
	}

	updateQuery := `
		UPDATE users SET points = points + $1 WHERE uuid = $2
	`
	if _, err := tx.ExecContext(ctx, updateQuery, taskEntity.Reward, userId); err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (u *User) UseReferrer(ctx context.Context, userId string, reward int) error {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return err
	}

	var userEntity = new(entities.User)
	refQuery := `
		SELECT uuid, name, points, referral FROM users WHERE uuid = $1
	`
	if err := tx.QueryRowContext(ctx, refQuery, userId).Scan(
		&userEntity.UserID, &userEntity.Name, &userEntity.Points, &userEntity.Referral,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.UserNotFoundError
		}
		tx.Rollback()
		return err
	}

	query := `
		UPDATE users 
		SET points = points + $1 
		WHERE uuid = $2
	`
	if _, err := tx.ExecContext(ctx, query, reward, userId); err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
