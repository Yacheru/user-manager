package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"user-manager/init/config"
	"user-manager/init/logger"
	"user-manager/internal/repository/mongodb"
	"user-manager/internal/repository/postgres"
	"user-manager/internal/server/http/router"
	"user-manager/pkg/constants"
)

type Server struct {
	http *http.Server
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	db, err := postgres.NewPostgresConnection(ctx, cfg.PostgresqlDSN)
	if err != nil {
		return nil, err
	}

	coll, err := mongodb.InitMongoDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	handler := initGin(cfg.ApiDebug)
	router.NewRouterAndComponents(handler.Group(cfg.ApiEntry), db, coll).Routes()

	return &Server{
		http: &http.Server{
			Addr:           fmt.Sprintf(":%d", cfg.ApiPort),
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}, nil
}

func (s *Server) Run() {
	go func() {
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error(), constants.ServerCategory)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func initGin(debug bool) *gin.Engine {
	var mode = gin.ReleaseMode
	if debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.New()

	engine.Use(gin.LoggerWithFormatter(logger.HTTPLogger))
	engine.Use(gin.Recovery())

	return engine
}
