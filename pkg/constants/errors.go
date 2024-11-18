package constants

import "errors"

var (
	CodeNotFoundError     = errors.New("code not found")
	CodeAlreadyExistError = errors.New("code already exist")

	NoDataError = errors.New("no data")

	TaskNotFoundError         = errors.New("task not found")
	TaskAlreadyCompletedError = errors.New("task already completed")

	UserNotFoundError     = errors.New("user not found")
	UserAlreadyExistError = errors.New("user already exist")
)
