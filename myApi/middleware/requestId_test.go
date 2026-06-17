package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestRequestIdMiddleware(t *testing.T) {

	okHandler := RequestIdMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	server := httptest.NewServer(okHandler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("request error: %v\n", err)
	}
	defer resp.Body.Close()

	id := resp.Header.Get(string("X-Request-ID"))

	if err := uuid.Validate(id); err != nil {
		t.Errorf("Request ID doesn`t initial : %v", err)
	} else {
		t.Logf("request id: %s", id)
	}
}
