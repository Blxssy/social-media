package app

import (
	"github.com/Blxssy/social-media/auth-service/internal/services/auth"
	"github.com/Blxssy/social-media/auth-service/internal/storage"
	"log/slog"

	grpcapp "github.com/Blxssy/social-media/auth-service/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    storage.Storage
}

func New(log *slog.Logger, grpcPort int, storage storage.Storage) *App {
	authService := auth.New(log, storage, storage)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
	}
}
