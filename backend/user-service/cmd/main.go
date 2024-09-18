package main

import (
	"log/slog"

	"github.com/Blxssy/social-media/user-service/internal/config"
	"github.com/Blxssy/social-media/user-service/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	log.Info("cfg", slog.Any("cfg", cfg))
}
