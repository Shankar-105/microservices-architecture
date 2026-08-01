package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONBodyRejectsUnknownFields(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/createorder",
		strings.NewReader(`{"user_id":"u-100","unexpected":true}`))
	recorder := httptest.NewRecorder()
	var body createOrderHTTPRequest

	if decodeJSONBody(recorder, request, &body) {
		t.Fatal("expected unknown fields to be rejected")
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestDecodeJSONBodyRejectsTrailingObject(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/createorder",
		strings.NewReader(`{"user_id":"u-100"}{"quantity":1}`))
	recorder := httptest.NewRecorder()
	var body createOrderHTTPRequest

	if decodeJSONBody(recorder, request, &body) {
		t.Fatal("expected trailing JSON object to be rejected")
	}
}

func TestValidIdentifier(t *testing.T) {
	for _, value := range []string{"", " ", strings.Repeat("x", 129)} {
		if validIdentifier(value) {
			t.Fatalf("expected identifier %q to be invalid", value)
		}
	}
	if !validIdentifier(" u-100 ") {
		t.Fatal("expected trimmed identifier to be valid")
	}
}
