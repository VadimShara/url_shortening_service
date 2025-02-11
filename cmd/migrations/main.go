package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/VadimShara/url_shortening_service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	PostgresCfg := config.Load().PGConfig

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		PostgresCfg.User,
		PostgresCfg.Password,
		PostgresCfg.Host,
		PostgresCfg.Port,
		PostgresCfg.DBName,
		PostgresCfg.SSLMode,
	)

	migrationsPath := config.Load().MigrationsPath

	migrateCmd := flag.String("migrate", "", "Specify 'up' or 'down' to run migrations")
	flag.Parse()

	m, err := migrate.New(
		migrationsPath,
		dbUrl,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch *migrateCmd {
	case "up":
		migrateUp(m)
	case "down":
		migrateDown(m)
	default:
		log.Fatal("--migrate flag not specified. Use '--migrate=up' or '--migrate=down'")
	}

	log.Println("Migrations applied successfully")
}

func migrateUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}

func migrateDown(m *migrate.Migrate) {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}
