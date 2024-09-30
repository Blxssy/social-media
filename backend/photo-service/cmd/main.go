package main

import (
	"github.com/Blxssy/social-media/photo-service/intenal/config"
	"github.com/Blxssy/social-media/photo-service/intenal/storage"
	"github.com/Blxssy/social-media/photo-service/pkg/logger"
	"log/slog"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.Env)
	logger.Info("cfg", slog.Any("cfg", cfg))

	store := storage.NewStorage(cfg, logger)

}
