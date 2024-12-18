package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-manager/internal/entities"
	"user-manager/pkg/constants"
)

// NewTask
//
//	@Summary	Create new task
//	@Security	JWT
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		input	body		entities.Task	true	"create new task"
//	@Success	200		{object}	handlers.Response
//	@Failure	400		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/admin/new/task [post]
func (h *Handler) NewTask(ctx *gin.Context) {
	var taskEntity = new(entities.Task)
	if err := ctx.ShouldBindJSON(taskEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body")
		return
	}

	task, err := h.s.AdminService.NewTask(ctx.Request.Context(), taskEntity)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "task has been successfully saved", task)
	return
}

// NewCode
//
//	@Summary	Create new referral code
//	@Security	JWT
//	@Tags		code
//	@Accept		json
//	@Produce	json
//	@Param		input	body		entities.Code	true	"create new referral code"
//	@Success	200		{object}	handlers.Response
//	@Failure	400		{object}	handlers.Response
//	@Failure	409		{object}	handlers.Response
//	@Failure	500		{object}	handlers.Response
//	@Router		/users/admin/new/code [post]
func (h *Handler) NewCode(ctx *gin.Context) {
	var codeEntity = new(entities.Code)
	if err := ctx.ShouldBindJSON(codeEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body")
		return
	}

	if codeEntity.Code == "" {
		NewErrorResponse(ctx, http.StatusBadRequest, "code cannot be empty")
		return
	}

	if err := h.s.AdminService.NewCode(ctx.Request.Context(), codeEntity); err != nil {
		if errors.Is(err, constants.CodeAlreadyExistError) {
			NewErrorResponse(ctx, http.StatusConflict, "code already exists")
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "code has been successfully saved", codeEntity.Code)
	return
}
