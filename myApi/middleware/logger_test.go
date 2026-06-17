package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {

	var out bytes.Buffer
	log.SetOutput(&out)
	defer log.SetOutput(os.Stderr)

	okHandler := Logger(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	server := httptest.NewServer(okHandler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Error of request: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code:\nexpexcted: %d\ngot: %d\n", http.StatusOK, resp.StatusCode)
	}

	//OK request check
	if !strings.Contains(out.String(), "200") {
		t.Errorf("logger message:\nexpected: %d\ngot: %s\n", http.StatusOK, out.String())
	}

	if !strings.Contains(out.String(), "GET") {
		t.Errorf("logger message:\nexpected: %s\ngot: %s\n", "GET", out.String())
	}

	if !strings.Contains(out.String(), "ms") {
		t.Errorf("logger message:\nexpected: %s\ngot: %s\n", "ms", out.String())
	}
}
