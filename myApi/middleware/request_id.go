package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const requestIdKey contextKey = "request-id"

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestId := uuid.New().String()

		ctx := context.WithValue(r.Context(), requestIdKey, requestId)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestId)

		next.ServeHTTP(w, r)
	})
}

func GetRequestId(ctx context.Context) string {
	if value, ok := ctx.Value(requestIdKey).(string); ok {
		return value
	}
	return "empty"
}
