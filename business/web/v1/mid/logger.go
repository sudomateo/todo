package mid

import (
	"context"
	"net/http"
	"time"
	"todo-api/web"

	"github.com/hashicorp/go-hclog"
)

func Logger(log hclog.Logger) web.Middleware {
	return func(h web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			err := h(ctx, w, r)

			log.Info("request completed",
				"trace_id", v.TraceID,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_address", r.RemoteAddr,
				"status_code", v.StatusCode,
				"duration", time.Since(v.Now),
			)

			return err
		}
	}
}
