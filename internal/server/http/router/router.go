package router

import (
	"github.com/gin-gonic/gin"
	"user-manager/internal/server/http/handlers"
)

type Router struct {
	handler *handlers.Handler
	router  *gin.RouterGroup
}

func NewRouterAndComponents(router *gin.RouterGroup) *Router {
	handler := handlers.NewHandler()

	return &Router{
		handler: handler,
		router:  router,
	}
}

func (r *Router) Routes() {
	{
		r.router.GET("/:id/status")
		r.router.GET("/leaderboard")
		r.router.POST("/:id/task/complete")
		r.router.POST("/:id/referrer")
	}
}
