package app

import (
	"go.uber.org/zap"
	grpcApp "sso/internal/app/grpc"
	"sso/internal/repository/sqlite"
	"sso/internal/service/auth"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *zap.Logger, storagePath string, grpcPort int, tokenTTL time.Duration) *App {

	repository, err := sqlite.NewRepository(storagePath)
	if err != nil {
		panic(err)
	}

	service := auth.NewAuthService(log, repository, repository, repository, tokenTTL)

	gRPCApp := grpcApp.NewApp(log, service, grpcPort)

	return &App{
		GRPCServer: gRPCApp,
	}
}
