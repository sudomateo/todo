package mid

import (
	"context"
	"errors"
	"net/http"
	"strings"
	v1 "todo-api/business/web/v1"
	"todo-api/web"
)

func Auth(authToken string) web.Middleware {
	return func(h web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authHeader := r.Header.Get("Authorization")
			fields := strings.Split(authHeader, " ")

			if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
				return v1.NewRequestError(errors.New("malformed authorization header"), http.StatusUnauthorized)
			}

			if fields[1] != authToken {
				return v1.NewRequestError(errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
			}

			return h(ctx, w, r)
		}
	}
}
