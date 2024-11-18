package repository

import (
	"context"
	"user-manager/internal/entities"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"

	"user-manager/internal/repository/mongodb"
	"user-manager/internal/repository/postgres"
)

type UserPostgres interface {
	NewUser(ctx context.Context, user *entities.NewUser) (*entities.User, error)
	GetStatus(ctx context.Context, id string) (*entities.User, error)
	Leaderboard(ctx context.Context, limit, offset string) (*[]entities.Leaderboard, error)
	TaskComplete(ctx context.Context, taskId, userId string) error
	UseReferrer(ctx context.Context, userId string, reward int) error
}

type TasksPostgres interface {
	NewTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
}

type CodesMongo interface {
	NewCode(ctx context.Context, code string, reward int) error
	FindCode(ctx context.Context, code string) (*entities.Code, error)
	RemoveCode(ctx context.Context, code string) error
}

type Repository struct {
	UserPostgres
	TasksPostgres
	CodesMongo
}

func NewRepository(db *sqlx.DB, coll *mongo.Collection) *Repository {
	return &Repository{
		UserPostgres:  postgres.NewUser(db),
		TasksPostgres: postgres.NewTasks(db),
		CodesMongo:    mongodb.NewCodes(coll),
	}
}
