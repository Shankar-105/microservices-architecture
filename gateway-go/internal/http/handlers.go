package http

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "gateway",
		"status":  "ok",
	})
}

func PlaceholderOrderCreateHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"message": "Phase 1 scaffold: route exists, implementation will be added in Phase 3",
	})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
