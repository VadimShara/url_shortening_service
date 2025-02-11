package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/VadimShara/url_shortening_service/internal/app"
	"github.com/VadimShara/url_shortening_service/internal/config"
	"github.com/VadimShara/url_shortening_service/pkg/logger"
)

func main() {
	cfg := config.Load()

	log := logger.SetupLogger(cfg.Env)

	application := app.New(cfg, log, cfg.GRPC.Port, cfg.MigrationsPath)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}
