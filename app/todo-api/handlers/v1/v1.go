package v1

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"path"
	"todo-api/app/todo-api/handlers/v1/todohandlers"
	"todo-api/business/todo"
	"todo-api/business/todo/store/tododb"
	"todo-api/business/web/v1/mid"
	"todo-api/web"

	"github.com/hashicorp/go-hclog"
)

type Config struct {
	Log       hclog.Logger
	DB        *sql.DB
	AuthToken string
	Version   string
}

func Routes(app *web.App, cfg Config) {
	const version = "/api/v1/"

	authMid := mid.Auth(cfg.AuthToken)

	todoHandler := todohandlers.Handler{
		Todo: todo.NewCore(tododb.NewStore(cfg.DB)),
	}

	app.Handle(http.MethodGet, path.Join(version, "/todos"), todoHandler.Query, authMid)
	app.Handle(http.MethodPost, path.Join(version, "/todos"), todoHandler.Create, authMid)
	app.Handle(http.MethodGet, path.Join(version, "/todos/:id"), todoHandler.QueryByID, authMid)
	app.Handle(http.MethodPatch, path.Join(version, "/todos/:id"), todoHandler.Update, authMid)
	app.Handle(http.MethodDelete, path.Join(version, "/todos/:id"), todoHandler.Delete, authMid)

	app.Handle(http.MethodGet, path.Join(version, "/version"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		data := struct {
			Version string `json:"version"`
		}{
			Version: cfg.Version,
		}

		return web.Respond(ctx, w, data, http.StatusOK)
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RequestError struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
