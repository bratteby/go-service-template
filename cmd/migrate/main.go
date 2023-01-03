package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func migrateUp(m *migrate.Migrate) {
	log.Println("***applying all up migrations***")

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("***no new migrations to apply***")
			os.Exit(0)
		}

		log.Fatal(fmt.Errorf("could not perform up migrations: %w", err))
	}
}

func migrateDown(m *migrate.Migrate) {
	log.Println("***applying all down migrations***")

	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("***no new migrations to apply***")
			os.Exit(0)
		}

		log.Fatal(fmt.Errorf("could not perform up migrations: %w", err))
	}
}

func getEnv(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func main() {
	var (
		host     = getEnv("POSTGRES_HOST", "localhost")
		port     = getEnv("POSTGRES_PORT", "5432")
		database = getEnv("POSTGRES_DB", "example")
		user     = getEnv("POSTGRES_USER", "postgres")
		password = getEnv("POSTGRES_PASSWORD", "")
		ssl      = getEnv("POSTGRES_SSL", "disable")

		migrationsFilePath = getEnv("MIGRATIONS_PATH", "migrations")
	)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, database, ssl,
	)

	m, err := migrate.New(
		"file://"+migrationsFilePath,
		connString,
	)
	if err != nil {
		log.Fatal(fmt.Errorf("could not create migrate instance: %w", err))
	}

	mode := flag.String("m", "up", "migrate up or down [up, down]")
	flag.Parse()

	switch *mode {
	case "up":
		migrateUp(m)
	case "down":
		migrateDown(m)
	default:
		log.Fatal(fmt.Errorf("invalid mode '%s'", *mode))
	}
}
