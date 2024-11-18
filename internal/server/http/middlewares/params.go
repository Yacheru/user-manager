package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-manager/internal/server/http/handlers"
)

func Params() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit := ctx.DefaultQuery("limit", "10")
		offset := ctx.DefaultQuery("offset", "0")

		intLimit, err := strconv.Atoi(limit)
		if err != nil || intLimit < 0 {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "limit must be a positive integer or zero")
			return
		}

		intOffset, err := strconv.Atoi(offset)
		if err != nil || intOffset < 0 {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "offset must be a positive integer or zero")
			return
		}

		ctx.Next()
	}
}
