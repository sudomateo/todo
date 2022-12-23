package handlers

import (
	"database/sql"
	"net/http"
	v1 "todo-api/app/todo-api/handlers/v1"
	"todo-api/business/web/v1/mid"
	"todo-api/web"

	"github.com/hashicorp/go-hclog"
)

type APIMuxConfig struct {
	Log       hclog.Logger
	DB        *sql.DB
	AuthToken string
	Version   string
}

func APIMux(cfg APIMuxConfig) http.Handler {
	app := web.NewApp(
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
	)

	v1.Routes(app, v1.Config{
		Log:       cfg.Log,
		DB:        cfg.DB,
		AuthToken: cfg.AuthToken,
		Version:   cfg.Version,
	})

	return app
}
