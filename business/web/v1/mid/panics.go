package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"todo-api/web"
)

func Panics() web.Middleware {
	return func(h web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if v := recover(); v != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", v, string(trace))
				}
			}()

			return h(ctx, w, r)
		}
	}
}
