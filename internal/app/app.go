package app

import (
	"log/slog"

	grpcapp "github.com/VadimShara/url_shortening_service/internal/app/grpc"
	"github.com/VadimShara/url_shortening_service/internal/config"
	"github.com/VadimShara/url_shortening_service/internal/repo"
	"github.com/VadimShara/url_shortening_service/internal/service"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	cfg config.Config,
	log *slog.Logger,
	grpcPort int,
	repoPath string,
) *App {
	repo := repo.NewRepository(cfg)

	urlService := service.New(log, repo, repo)

	log.Info(cfg.Storage)

	grpcApp := grpcapp.New(log, urlService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
