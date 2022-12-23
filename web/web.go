package web

import (
	"context"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	mux *httptreemux.ContextMux
	mw  []Middleware
}

func NewApp(mw ...Middleware) *App {
	return &App{
		mux: httptreemux.NewContextMux(),
		mw:  mw,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, ctxKey{}, &v)

		if err := handler(ctx, w, r); err != nil {
			// TODO
		}
	}

	a.mux.Handle(method, path, h)
}
