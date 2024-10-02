package main

import (
	"github.com/Blxssy/social-media/photo-service/intenal/app"
	"github.com/Blxssy/social-media/photo-service/intenal/config"
	"github.com/Blxssy/social-media/photo-service/intenal/storage"
	"github.com/Blxssy/social-media/photo-service/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.Env)
	logger.Info("cfg", slog.Any("cfg", cfg))

	store := storage.NewStorage(cfg, logger)

	application := app.New(logger, cfg.GRPC.Port, store)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("Server stopped")
}
