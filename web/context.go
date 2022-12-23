package web

import (
	"context"
	"time"
)

type ctxKey struct{}

type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(ctxKey{}).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey{}).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(ctxKey{}).(*Values)
	if !ok {
		return time.Now()
	}
	return v.Now
}

func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(ctxKey{}).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}
