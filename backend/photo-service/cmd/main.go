package main

import (
	"github.com/Blxssy/social-media/photo-service/intenal/config"
	"github.com/Blxssy/social-media/photo-service/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.Env)
	_ = logger
}
