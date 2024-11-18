package service

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"sync"
	"user-manager/internal/jwt"

	"user-manager/internal/entities"
	"user-manager/internal/repository"
	"user-manager/pkg/constants"
)

type User struct {
	userPostgres repository.UserPostgres
	codesMongo   repository.CodesMongo

	mu sync.Mutex
}

func NewUser(userPostgres repository.UserPostgres, codesMongo repository.CodesMongo) *User {
	return &User{
		userPostgres: userPostgres,
		codesMongo:   codesMongo,
		mu:           sync.Mutex{},
	}
}

func (u *User) NewUser(ctx context.Context, user *entities.NewUser) (string, error) {
	id := uuid.NewString()
	user.UserID = &id

	_, err := u.userPostgres.NewUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return "", constants.UserAlreadyExistError
		}
		return "", err
	}

	token, err := jwt.NewToken(*user.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *User) GetStatus(ctx context.Context, id string) (*entities.User, error) {
	user, err := u.userPostgres.GetStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Leaderboard(ctx context.Context, limit, offset string) (*[]entities.Leaderboard, error) {
	users, err := u.userPostgres.Leaderboard(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) TaskComplete(ctx context.Context, taskId, userId string) error {
	if err := u.userPostgres.TaskComplete(ctx, taskId, userId); err != nil {
		return err
	}
	return nil
}

func (u *User) UseReferrer(ctx context.Context, userId, code string) (int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	codeEntity, err := u.codesMongo.FindCode(ctx, code)
	if err != nil {
		return 0, err
	}

	if err := u.userPostgres.UseReferrer(ctx, userId, codeEntity.Reward); err != nil {
		return 0, err
	}

	if err := u.codesMongo.RemoveCode(ctx, code); err != nil {
		return 0, err
	}

	return codeEntity.Reward, nil
}
