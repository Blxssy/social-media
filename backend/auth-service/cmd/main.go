package main

import (
	"github.com/Blxssy/social-media/auth-service/internal/app"
	"github.com/Blxssy/social-media/auth-service/internal/config"
	"github.com/Blxssy/social-media/auth-service/internal/storage"
	"github.com/Blxssy/social-media/auth-service/pkg/logger"
	"github.com/Blxssy/social-media/auth-service/pkg/token"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	godotenv.Load()
	token.InitJWTKey()

	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.Env)

	store := storage.NewStorage(logger, cfg)

	application := app.New(logger, cfg.GRPC.Port, store)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("Gracefully stopped")
}
