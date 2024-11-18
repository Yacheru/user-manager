package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"user-manager/internal/entities"
	"user-manager/internal/repository"
	"user-manager/pkg/constants"
)

type Admin struct {
	codesMongo    repository.CodesMongo
	tasksPostgres repository.TasksPostgres
}

func (c *Admin) NewTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	id := uuid.NewString()
	task.TaskID = &id

	task, err := c.tasksPostgres.NewTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func NewAdmin(codesMongo repository.CodesMongo, tasksPostgres repository.TasksPostgres) *Admin {
	return &Admin{codesMongo: codesMongo, tasksPostgres: tasksPostgres}
}

func (c *Admin) NewCode(ctx context.Context, code *entities.Code) error {
	savedCode, err := c.codesMongo.FindCode(ctx, code.Code)
	if err != nil {
		if !errors.Is(err, constants.CodeNotFoundError) {
			return err
		}
		if err := c.codesMongo.NewCode(ctx, code.Code, code.Reward); err != nil {
			return err
		} else {
			return nil
		}
	}

	if savedCode != nil {
		return constants.CodeAlreadyExistError
	}

	return nil
}
