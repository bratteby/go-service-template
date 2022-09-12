package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bratteby/go-service-template/example"
	"github.com/bratteby/go-service-template/httpserver"
	"github.com/bratteby/go-service-template/logging"
	"github.com/bratteby/go-service-template/postgres"
)

func main() {
	// Environment variables.
	var (
		ADDRESS = getEnv("HTTP_ADDRESS", ":80")

		POSTGRES_HOST     = getEnv("POSTGRES_HOST", "localhost")
		POSTGRES_PORT     = getEnv("POSTGRES_PORT", "5432")
		POSTGRES_DATABASE = getEnv("POSTGRES_DB", "example")
		POSTGRES_USER     = getEnv("POSTGRES_USER", "postgres")
		POSTGRES_PASSWORD = getEnv("POSTGRES_PASSWORD", "")
		POSTGRES_SSL      = getEnv("POSTGRES_SSL", "disable")
	)

	errorChannel := make(chan error)

	logger := logging.New(nil, logging.Config{
		Level:         logging.InfoLevel,
		WithTimeStamp: true,
		Options:       []logging.Option{},
	})

	// Repositories
	dbPool, err := postgres.NewPool(postgres.ConnectionConfig{
		User:     POSTGRES_USER,
		Password: POSTGRES_PASSWORD,
		Host:     POSTGRES_HOST,
		Port:     POSTGRES_PORT,
		DB:       POSTGRES_DATABASE,
		SSL:      POSTGRES_SSL,
	})
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	exampleRepository := &postgres.ExampleRepository{
		DB: dbPool,
	}

	// Services.
	exampleService := example.Service{
		ExampleRepository: exampleRepository,
		Logger:            logger,
	}

	// HTTP.
	go func() {
		httpServer := httpserver.Server{
			Address:        ADDRESS,
			ExampleService: exampleService,
			Logger:         logger,
		}

		logger.Infof("starting server on: '%s'", ADDRESS)

		errorChannel <- httpServer.Start()
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errorChannel <- fmt.Errorf("got signal: %s", <-c)
	}()

	if err := <-errorChannel; err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func getEnv(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}
