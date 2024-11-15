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

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()

	cfg := &config.CFG

	if err := config.InitConfig(); err != nil {
		cancel()
	}

	log := logger.InitLogger(cfg.ApiDebug)

	httpServer, err := server.NewServer(cfg, log)
	if err != nil {
		cancel()
	}

	logger.Info("service started", constants.MainCategory)

	<-ctx.Done()

	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error(err.Error(), constants.MainCategory)
		}
		logger.Info("server stopped", constants.MainCategory)
	}

	logger.Info("service exited", constants.MainCategory)
}
