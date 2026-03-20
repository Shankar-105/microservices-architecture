package http

import (
	"net/http"

	"github.com/shank/bookstore-microservices/gateway-go/internal/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthHandler)
	mux.HandleFunc("/v1/orders", PlaceholderOrderCreateHandler)

	return middleware.RequestID(mux)
}
