package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/app/todo-api/handlers"
	"todo-api/business/database"

	"github.com/hashicorp/go-hclog"
	_ "github.com/lib/pq"
)

var (
	defaultLogLevel = "info"
	defaultAddress  = ":7836"
	defaultVersion= "development"
)

type config struct {
	AuthToken   string
	DatabaseURL string

	Address  string
	LogLevel string
	Version  string
}

func main() {
	log := hclog.New(&hclog.LoggerOptions{
		Name:  "todo-api",
		Level: hclog.LevelFromString(defaultLogLevel),
	})

	if err := run(log); err != nil {
		log.Error("startup", "error", err)
		os.Exit(1)
	}
}

func run(log hclog.Logger) error {
	config, err := configFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	log.SetLevel(hclog.LevelFromString(config.LogLevel))

	log.Info("starting service", "version", config.Version)
	defer log.Info("shutdown complete")

	u, err := url.Parse(config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("invalid database url: %w", err)
	}

	log.Info("startup", "status", "initializing database", "host", u.Host)
	db, err := database.Open(config.DatabaseURL)
	if err != nil {
		return err
	}

	log.Info("startup", "status", "running database migrations", "host", u.Host)
	migrateCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := database.Migrate(migrateCtx, db); err != nil {
		return err
	}

	log.Info("startup", "status", "initializing api server")

	mux := handlers.APIMux(handlers.APIMuxConfig{
		Log:       log,
		DB:        db,
		AuthToken: config.AuthToken,
		Version:   config.Version,
	})

	api := http.Server{
		Addr:         config.Address,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	errCh := make(chan error, 1)

	go func() {
		log.Info("startup", "status", "api server started", "address", api.Addr)
		errCh <- api.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("api error: %w", err)

	case sig := <-shutdown:
		log.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop api server gracefully: %w", err)
		}
	}

	return nil
}

func configFromEnv() (config, error) {
	databaseURL := os.Getenv("TODO_DATABASE_URL")
	if databaseURL == "" {
		return config{}, errors.New("TODO_DATABASE_URL is required")
	}

	authToken := os.Getenv("TODO_AUTH_TOKEN")
	if authToken == "" {
		return config{}, errors.New("TODO_AUTH_TOKEN is required")
	}

	address := os.Getenv("TODO_ADDR")
	if address == "" {
		address = defaultAddress
	}

	version := os.Getenv("TODO_VERSION")
	if version == "" {
		version = defaultVersion
	}

	logLevel := os.Getenv("TODO_LOG_LEVEL")
	if logLevel == "" {
		logLevel = defaultLogLevel
	}
	if hclog.LevelFromString(logLevel) == hclog.NoLevel {
		return config{}, fmt.Errorf("invalid log level %s", logLevel)
	}

	cfg := config{
		DatabaseURL: databaseURL,
		AuthToken:   authToken,
		Address:     address,
		Version:     version,
		LogLevel:    logLevel,
	}

	return cfg, nil
}
