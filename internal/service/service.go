package service

import (
	"context"
	"user-manager/internal/entities"
	"user-manager/internal/repository"
)

type UserService interface {
	NewUser(ctx context.Context, user *entities.NewUser) (string, error)
	GetStatus(ctx context.Context, id string) (*entities.User, error)
	Leaderboard(ctx context.Context, limit, offset string) (*[]entities.Leaderboard, error)
	TaskComplete(ctx context.Context, taskId, userId string) error
	UseReferrer(ctx context.Context, userId, code string) (int, error)
}

type AdminService interface {
	NewCode(ctx context.Context, code *entities.Code) error
	NewTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
}

type Service struct {
	UserService
	AdminService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		UserService:  NewUser(repo.UserPostgres, repo.CodesMongo),
		AdminService: NewAdmin(repo.CodesMongo, repo.TasksPostgres),
	}
}
