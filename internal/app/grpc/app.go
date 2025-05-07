package grpcApp

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	authgRPC "sso/internal/grpc/auth"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *zap.Logger, authService authgRPC.AuthService, port int) *App {
	gRPCServer := grpc.NewServer()

	authgRPC.Register(gRPCServer, authService)

	return &App{log: log, port: port, gRPCServer: gRPCServer}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	a.log.Info("starting gRPC server", zap.String("address", l.Addr().String()))

	if err = a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
