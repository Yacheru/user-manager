package middlewares

import (
	"github.com/gin-gonic/gin"

	"user-manager/internal/jwt"
	"user-manager/internal/server/http/handlers"

	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("authorization")
		if auth == "" {
			handlers.NewErrorResponse(ctx, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		err := jwt.ValidateToken(parts[1])
		if err != nil {
			handlers.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx.Next()
	}
}
