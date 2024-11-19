package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-manager/internal/entities"
	"user-manager/internal/service"
	"user-manager/pkg/constants"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

// NewUser
//
//	@Summary	Create new user
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		input	body		entities.NewUser	true "Name for new user"
//	@Success	200		{object}	handlers.Response
//	@Failure	400		{object}	handlers.Response
//	@Failure	409		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/new [post]
func (h *Handler) NewUser(ctx *gin.Context) {
	var userEntity = new(entities.NewUser)
	if err := ctx.ShouldBindJSON(&userEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body")
		return
	}

	token, err := h.s.UserService.NewUser(ctx.Request.Context(), userEntity)
	if err != nil {
		if errors.Is(err, constants.UserAlreadyExistError) {
			NewErrorResponse(ctx, http.StatusConflict, "user with that name already exists")
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "new user", token)
	return
}

// GetStatus
//
//	@Summary	Get user status by user ID
//	@Security	JWT
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		input	path		string	true	"User ID (uuid v4)"
//	@Success	200		{object}	handlers.Response
//	@Failure	404		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/:id/status [get]
func (h *Handler) GetStatus(ctx *gin.Context) {
	user, err := h.s.GetStatus(ctx, ctx.Param("id"))
	if err != nil {
		if errors.Is(err, constants.UserNotFoundError) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "status", user)
	return
}

// Leaderboard
//
//	@Summary	Get a leaderboard of the richest users
//	@Security	JWT
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		string	false	"Default: 10"
//	@Param		offset	query		string	false	"Default: 0"
//	@Success	200		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/leaderboard [get]
func (h *Handler) Leaderboard(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "10")
	offset := ctx.DefaultQuery("offset", "0")

	users, err := h.s.UserService.Leaderboard(ctx.Request.Context(), limit, offset)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "leaderboard", users)
	return
}

// TaskComplete
//
//	@Summary	Complete some task
//	@Security	JWT
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		task_id	query		string	true	"Task ID (uuid v4)"
//	@Param		id		path		string	true	"User ID (uuid v4)"
//	@Success	200		{object}	handlers.Response
//	@Failure	404		{object}	handlers.Response
//	@Failure	409		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/:id/task/complete [post]
func (h *Handler) TaskComplete(ctx *gin.Context) {
	taskId := ctx.Query("task_id")
	userId := ctx.Param("id")

	if err := h.s.TaskComplete(ctx, taskId, userId); err != nil {
		if errors.Is(err, constants.TaskNotFoundError) || errors.Is(err, constants.UserNotFoundError) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, constants.TaskAlreadyCompletedError) {
			NewErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "task completed successfully", nil)
	return
}

// UseReferrer
//
//	@Summary	Use referral code
//	@Security	JWT
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		code	query		string	true	"Referral code"
//	@Param		id		path		string	true	"User ID (uuid v4)"
//	@Success	200		{object}	handlers.Response
//	@Failure	404		{object}	handlers.Response
//	@Failure	400		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/:id/referrer [post]
func (h *Handler) UseReferrer(ctx *gin.Context) {
	code := ctx.Query("code")
	userId := ctx.Param("id")

	if code == "" {
		NewErrorResponse(ctx, http.StatusBadRequest, "code is required")
		return
	}

	reward, err := h.s.UserService.UseReferrer(ctx.Request.Context(), userId, code)
	if err != nil {
		if errors.Is(err, constants.UserNotFoundError) || errors.Is(err, constants.CodeNotFoundError) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "your reward", reward)
	return
}
