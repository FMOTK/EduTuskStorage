package middleware

import (
	"log"
	"net/http"
	"time"
)

type requestStatusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *requestStatusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *requestStatusRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		respWriter := &requestStatusRecorder{
			ResponseWriter: w,
			statusCode:     200,
		}

		start := time.Now

		next.ServeHTTP(respWriter, r)

		rDuration := time.Since(start())
		rID := GetRequestId(r.Context())

		log.Printf("[%s] %s %s %d %s",
			rID,
			r.Method,
			r.URL.String(),
			respWriter.statusCode,
			rDuration.String(),
		)
	})
}
