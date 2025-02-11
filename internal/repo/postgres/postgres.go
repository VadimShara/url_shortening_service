package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	errs "github.com/VadimShara/url_shortening_service/pkg/errs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// func PGConnectionStr(config config.PostresConfig) string {
// 	return fmt.Sprintf(
// 		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
// 		config.User,
// 		config.Password,
// 		config.Host,
// 		config.Port,
// 		config.DBName,
// 		config.SSLMode,
// 	)
// }

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(user, password, host string, port int, dbName, sslMode string) *DB {
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("unable to parse db url")
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("unable to connect to db")
	}

	return &DB{Pool: pool}
}

func (db *DB) CLose() {
	db.Pool.Close()
}

func (d *DB) SaveUrl(ctx context.Context, urlToSave, alias string) (string, error) {
	const op = "repo.postgres.SaveUrl"

	const query = "INSERT INTO url (url, alias) VALUES ($1, $2)"

	_, err := d.Pool.Exec(ctx, query, urlToSave, alias)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "url_url_key":
				var existingAlias string
				err = d.Pool.QueryRow(ctx,
					"SELECT alias FROM url WHERE url = $1", urlToSave).Scan(&existingAlias)
				if err != nil {
					return "", fmt.Errorf("%s: %w", op, err)
				}
				return existingAlias, errs.ErrUrlExists
			case "url_alias_key":
				return "", fmt.Errorf("%s: %w", op, errs.ErrAliasExists)
			}
		}
		return "", fmt.Errorf("%s : failed to save url and alias: %w", op, err)
	}

	return alias, nil
}

func (d *DB) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "repo.postgres.GetUrl"

	const query = "SELECT url FROM url WHERE alias = $1"

	var url string

	err := d.Pool.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, errs.ErrUrlNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
