package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"user-manager/internal/server/http/handlers"
)

func TaskID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("task_id")
		if id == "" {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "No task id supplied")
			return
		}

		if _, err := uuid.Parse(id); err != nil {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "Invalid task id supplied")
			return
		}

		ctx.Next()
	}
}
