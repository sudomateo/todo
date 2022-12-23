package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/sudomateo/todo/database"
	"github.com/sudomateo/todo/todo"
	"github.com/sudomateo/todo/todo/stores/tododb"
	"github.com/sudomateo/todo/todo/stores/todomemory"
)

//go:embed views
var viewsFS embed.FS

//go:embed public
var publicFS embed.FS

const (
	defaultAddress  = ":8080"
	defaultLogLevel = "info"
	defaultVersion  = "1.0.0"
)

func main() {
	log := hclog.New(&hclog.LoggerOptions{
		Name:  "todo",
		Level: hclog.LevelFromString(defaultLogLevel),
	})

	if err := run(log); err != nil {
		log.Error("startup", "error", err)
		os.Exit(1)
	}
}

func run(log hclog.Logger) error {
	cfg, err := configFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	log.SetLevel(hclog.LevelFromString(cfg.LogLevel))

	todoCore := todo.NewCore(todomemory.NewStore())

	if cfg.Database.Host != "" {
		databaseURL := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(cfg.Database.User, cfg.Database.Password),
			Host:     cfg.Database.Host,
			Path:     cfg.Database.Name,
			RawQuery: cfg.Database.Parameters,
		}

		log.Info("startup", "status", "initializing database", "host", databaseURL.Host)
		db, err := database.Open(databaseURL.String())
		if err != nil {
			return fmt.Errorf("could not open database: %w", err)
		}

		log.Info("startup", "status", "running database migrations", "host", databaseURL.Host)
		migrateCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := database.Migrate(migrateCtx, db); err != nil {
			return fmt.Errorf("could not migrate database: %w", err)
		}

		todoCore = todo.NewCore(tododb.NewStore(db))
	}

	log.Info("starting service", "version", cfg.Version)
	defer log.Info("shutdown complete")

	a := App{
		Log:      log,
		TodoCore: todoCore,
		Version:  cfg.Version,
	}

	e := echo.New()
	e.StaticFS("static", echo.MustSubFS(publicFS, "public"))
	e.Renderer = &Template{
		template: template.Must(template.ParseFS(viewsFS, "views/*.tmpl")),
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				log.Error("error serving request", "error", err)
				return err
			}

			return nil
		}
	})
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			now := time.Now()

			defer func() {
				if !strings.HasPrefix(c.Request().URL.Path, "/static") {
					log.Info("request completed",
						"method", c.Request().Method,
						"path", c.Request().URL.Path,
						"remoteaddr", c.Request().RemoteAddr,
						"statuscode", c.Response().Status,
						"since", time.Since(now),
					)
				}
			}()

			return next(c)
		}
	})

	e.GET("/", a.Root)
	e.GET("/api/todo", a.Query)
	e.GET("/api/todo/:id", a.QueryByID)
	e.POST("/api/todo", a.Create)
	e.PATCH("/api/todo/:id", a.Update)
	e.DELETE("/api/todo/:id", a.Delete)

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      e,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	serverErrors := make(chan error, 1)

	go func() {
		log.Info("startup", "status", "server started", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}

// App represents our web application.
type App struct {
	Log      hclog.Logger
	TodoCore *todo.Core
	Version  string
}

// Root serves the web application.
func (a *App) Root(c echo.Context) error {
	todos, err := a.TodoCore.Query(c.Request().Context())
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	data := struct {
		Todos   []todo.Todo
		Version string
	}{
		Todos:   todos,
		Version: a.Version,
	}

	return c.Render(http.StatusOK, "index.html.tmpl", data)
}

// Query fetches all todos.
func (a *App) Query(c echo.Context) error {
	todos, err := a.TodoCore.Query(c.Request().Context())
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return c.NoContent(http.StatusNotFound)
		default:
			return fmt.Errorf("query: %w", err)
		}
	}

	return c.JSON(http.StatusOK, todos)
}

// QueryByID fetches a single todo by its ID.
func (a *App) QueryByID(c echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id format")
	}

	t, err := a.TodoCore.QueryByID(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return c.NoContent(http.StatusNotFound)
		default:
			return fmt.Errorf("query: %w", err)
		}
	}

	return c.JSON(http.StatusOK, t)
}

// Create creates a todo.
func (a *App) Create(c echo.Context) error {
	var params todo.TodoCreateParams

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	t, err := a.TodoCore.Create(c.Request().Context(), params)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	return c.JSON(http.StatusCreated, t)
}

// Update updates a todo.
func (a *App) Update(c echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id format")
	}

	t, err := a.TodoCore.QueryByID(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return c.NoContent(http.StatusNotFound)
		default:
			return fmt.Errorf("query by id [%s]: %w", id, err)
		}
	}

	var params todo.TodoUpdateParams

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	t, err = a.TodoCore.Update(c.Request().Context(), t, params)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return c.JSON(http.StatusOK, t)
}

// Delete deletes a todo.
func (a *App) Delete(c echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id format")
	}

	t, err := a.TodoCore.QueryByID(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return c.NoContent(http.StatusNoContent)
		default:
			return fmt.Errorf("query by id [%s]: %w", id, err)
		}
	}

	if err := a.TodoCore.Delete(c.Request().Context(), t); err != nil {
		return fmt.Errorf("delete [%s]: %w", id, err)
	}

	return c.NoContent(http.StatusNoContent)
}

type Template struct {
	template *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

// Config represents the application configuration.
type Config struct {
	Address  string
	Database Database
	LogLevel string
	Version  string
}

type Database struct {
	User       string
	Password   string
	Host       string
	Name       string
	Parameters string
}

// configFromEnv reads application configuration from the environment.
func configFromEnv() (Config, error) {
	database := Database{
		User:       os.Getenv("TODO_DATABASE_USER"),
		Password:   os.Getenv("TODO_DATABASE_PASSWORD"),
		Host:       os.Getenv("TODO_DATABASE_HOST"),
		Name:       os.Getenv("TODO_DATABASE_NAME"),
		Parameters: os.Getenv("TODO_DATABASE_PARAMETERS"),
	}

	if database.Parameters == "" {
		database.Parameters = "sslmode=disable"
	}

	address := os.Getenv("TODO_ADDR")
	if address == "" {
		address = defaultAddress
	}

	logLevel := os.Getenv("TODO_LOG_LEVEL")
	if logLevel == "" {
		logLevel = defaultLogLevel
	}
	if hclog.LevelFromString(logLevel) == hclog.NoLevel {
		return Config{}, fmt.Errorf("invalid log level %s", logLevel)
	}

	version := os.Getenv("TODO_VERSION")
	if version == "" {
		version = defaultVersion
	}

	cfg := Config{
		Database: database,
		Address:  address,
		Version:  version,
		LogLevel: logLevel,
	}

	return cfg, nil
}
