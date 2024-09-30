package app

import (
	grpcapp "github.com/Blxssy/social-media/photo-service/intenal/app/grpc"
	"github.com/Blxssy/social-media/photo-service/intenal/services/photo"
	"github.com/Blxssy/social-media/photo-service/intenal/storage"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    storage.Storage
}

func New(log *slog.Logger, grpcPort int, storage storage.Storage) *App {
	photoService := photo.New(log, storage)
	grpcApp := grpcapp.New(log, photoService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
	}
}
