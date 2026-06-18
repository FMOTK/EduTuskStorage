package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	callCount := 0
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			panic("test panic")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := httptest.NewServer(RecoveryMiddleware(panicHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error of get request: %v", err)
	}

	//status code check
	expectedStatus := http.StatusInternalServerError
	actualStatus := resp.StatusCode
	if expectedStatus != actualStatus {
		t.Errorf("status code:\nexpected: %d\ngot: %d\n", expectedStatus, actualStatus)
	}

	//header check
	expectedHeader := "text/plain; charset=utf-8"
	actualHeader := resp.Header.Get("Content-Type")
	if expectedHeader != actualHeader {
		t.Errorf("header:\nexpected: %s\ngot: %s\n", expectedHeader, actualHeader)
	}

	//body check
	expectedBody := "Internal server Error\n"
	actualBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error of reading response body: %v\n", err)
	}
	defer resp.Body.Close()

	if expectedBody != string(actualBody) {
		t.Errorf("body:\nexpected: %s\ngot: %s\n", expectedBody, actualBody)
	}

	resp2, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Second request don`t processing. Error: %v\n", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Second request status code:\nexpexted: %d\n got: %d\n", http.StatusOK, resp2.StatusCode)
	}
}
