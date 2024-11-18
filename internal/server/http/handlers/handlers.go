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
