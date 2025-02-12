package repo

import (
	"context"

	"github.com/VadimShara/url_shortening_service/internal/config"
	"github.com/VadimShara/url_shortening_service/internal/repo/inMemory"
	"github.com/VadimShara/url_shortening_service/internal/repo/postgres"
)

type Repository interface {
	SaveUrl(ctx context.Context, urlToSave, alias string) (string, error)
	GetUrl(ctx context.Context, alias string) (string, error)
}

func NewRepository(cfg config.Config) Repository {
	repo := cfg.Storage

	switch repo {
	case "postgres":
		return postgres.NewDB(cfg.PGConfig)
	default:
		return inMemory.NewDB()
	}
}
