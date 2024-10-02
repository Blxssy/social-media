package main

import (
	"github.com/Blxssy/social-media/auth-service/internal/app"
	"github.com/Blxssy/social-media/auth-service/internal/config"
	"github.com/Blxssy/social-media/auth-service/internal/storage"
	"github.com/Blxssy/social-media/auth-service/pkg/logger"
	"github.com/Blxssy/social-media/auth-service/pkg/token"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return 
	}
	token.InitJWTKey()

	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.Env)

	logger.Info("cfg", slog.Any("cfg", cfg))

	store := storage.NewStorage(logger, cfg)

	application := app.New(logger, cfg.GRPC.Port, store)

	go func() {
		application.GRPCServer.MustRun()
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Prometheus metrics available at /metrics on port 9090")
	go func() {
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("Gracefully stopped")
}
