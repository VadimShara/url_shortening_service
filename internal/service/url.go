package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/VadimShara/url_shortening_service/internal/lib/generateAlias"
	"github.com/VadimShara/url_shortening_service/pkg/errs"
	"github.com/VadimShara/url_shortening_service/pkg/logger"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 10

type SaveUrlRequest struct {
	Url string `validate:"required,url"`
}

type UrlSaver interface {
	SaveUrl(
		ctx context.Context,
		urlToSave,
		alias string,
	) (string, error)
}

type UrlRedirecter interface {
	GetUrl(
		ctx context.Context,
		alias string,
	) (string, error)
}

type Url struct {
	log           *slog.Logger
	urlSaver      UrlSaver
	urlRedirecter UrlRedirecter
}

func New(
	log *slog.Logger,
	urlSaver UrlSaver,
	urlRedirecter UrlRedirecter,
) *Url {
	return &Url{
		urlSaver:      urlSaver,
		urlRedirecter: urlRedirecter,
		log:           log,
	}
}

func (u *Url) SaveUrl(ctx context.Context, url string) (string, error) {
	const op = "Url.SaveUrl" //operation

	validate := validator.New()

	req := SaveUrlRequest{
		Url: url,
	}

	err := validate.Struct(req)
	if err != nil {
		u.log.Error("invalid request", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	alias := generateAlias.NewAlias(aliasLength)

	alias, err = u.urlSaver.SaveUrl(ctx, url, alias)
	if err != nil {
		if errors.Is(err, errs.ErrUrlExists) {
			u.log.Info("url already exist")

			return alias, nil
		} else if errors.Is(err, errs.ErrAliasExists) {
			return u.SaveUrl(ctx, url)
		} else {
			u.log.Error("failed to save url", logger.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return alias, nil
}

func (u *Url) RedirectUrl(ctx context.Context, alias string) (string, error) {
	const op = "Url.RedirectUrl" //operation

	url, err := u.urlRedirecter.GetUrl(ctx, alias)
	if err != nil {
		if errors.Is(err, errs.ErrUrlNotFound) {
			u.log.Info("url not found", "alias", alias)

			return "", fmt.Errorf("%s: %w", op, errs.ErrUrlNotFound)
		}

		u.log.Error("failed to get url", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
