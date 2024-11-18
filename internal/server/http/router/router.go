package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"user-manager/internal/repository"
	"user-manager/internal/server/http/handlers"
	"user-manager/internal/server/http/middlewares"
	"user-manager/internal/service"
)

type Router struct {
	handler *handlers.Handler
	router  *gin.RouterGroup
}

func NewRouterAndComponents(router *gin.RouterGroup, db *sqlx.DB, coll *mongo.Collection) *Router {
	repo := repository.NewRepository(db, coll)
	srv := service.NewService(repo)
	handler := handlers.NewHandler(srv)

	return &Router{
		handler: handler,
		router:  router,
	}
}

func (r *Router) Routes() {
	r.router.POST("/new", r.handler.NewUser)

	{
		r.router.GET("/leaderboard", middlewares.Auth(), middlewares.Params(), r.handler.Leaderboard)

		r.router.GET("/:id/status", middlewares.Auth(), middlewares.ID(), r.handler.GetStatus)
		r.router.POST("/:id/task/complete", middlewares.Auth(), middlewares.ID(), middlewares.TaskID(), r.handler.TaskComplete)
		r.router.POST("/:id/referrer", middlewares.Auth(), middlewares.ID(), r.handler.UseReferrer)
	}

	admin := r.router.Group("/admin", middlewares.Auth())
	{
		admin.POST("/new/task", r.handler.NewTask)
		admin.POST("/new/code", r.handler.NewCode)
	}
}
