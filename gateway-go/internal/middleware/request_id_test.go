package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDAcceptsSafeHeader(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Context().Value(RequestIDKey)
		if value != "trace-123" {
			t.Fatalf("unexpected request ID context: %v", value)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set("X-Request-ID", "trace-123")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)
	if recorder.Header().Get("X-Request-ID") != "trace-123" {
		t.Fatal("expected safe request ID to be preserved")
	}
}

func TestRequestIDReplacesUnsafeHeader(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set("X-Request-ID", strings.Repeat("x", 129))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)
	value := recorder.Header().Get("X-Request-ID")
	if value == "" || len(value) > 128 {
		t.Fatalf("unexpected generated request ID: %q", value)
	}
}
