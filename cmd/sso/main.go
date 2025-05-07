package main

import (
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/pkg/logging"
	"syscall"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.GetLogger()

	application := app.NewApp(logger, cfg.StoragePath, cfg.GRPC.Port, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("application stopped")
}
