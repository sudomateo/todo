package mid

import (
	"context"
	"net/http"
	"todo-api/business"
	v1 "todo-api/business/web/v1"
	"todo-api/web"

	"github.com/hashicorp/go-hclog"
)

func Errors(log hclog.Logger) web.Middleware {
	return func(h web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := h(ctx, w, r); err != nil {
				log.Error("request error", "trace_id", web.GetTraceID(ctx), "error", err)

				var errResp v1.ErrorResponse
				var status int
				switch {
				case v1.IsRequestError(err):
					reqErr := v1.GetRequestError(err)
					errResp = v1.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case business.IsValidationError(err):
					vErr := business.GetValidationError(err)
					errResp = v1.ErrorResponse{
						Error: vErr.Error(),
					}
					status = http.StatusBadRequest

				default:
					errResp = v1.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, errResp, status); err != nil {
					return err
				}
			}

			return nil
		}
	}
}
