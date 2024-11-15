package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"user-manager/init/config"
	"user-manager/init/logger"
	"user-manager/pkg/constants"
)

type Server struct {
	http *http.Server
}

func NewServer(cfg *config.Config, log *logrus.Logger) (*Server, error) {
	handler := initGin(cfg.ApiDebug)

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
		if err := s.http.ListenAndServe(); err != nil {
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
