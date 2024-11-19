package main

import (
	"context"
	"os/signal"
	"syscall"

	"user-manager/init/config"
	"user-manager/init/logger"
	"user-manager/internal/server"
	"user-manager/pkg/constants"
)

//	@title		User-Manager
//	@version	1.0

//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@host						localhost
//	@BasePath					/

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()

	cfg := &config.CFG

	if err := config.InitConfig(); err != nil {
		cancel()
	}

	_ = logger.InitLogger(cfg.ApiDebug)

	httpServer, err := server.NewServer(ctx, cfg)
	if err != nil {
		cancel()
	}

	if httpServer != nil {
		httpServer.Run()
	}

	logger.InfoF("service started on %d port", constants.MainCategory, cfg.ApiPort)

	<-ctx.Done()

	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error(err.Error(), constants.MainCategory)
		}
		logger.Info("http-server stopped", constants.MainCategory)
	}

	logger.Info("service exited", constants.MainCategory)
}
