package http

import (
	"net/http"

	"microservices/gateway-go/internal/clients"
	"microservices/gateway-go/internal/middleware"
)

func NewRouter(orderClient *clients.OrderClient) http.Handler {
	handler := NewHandler(orderClient)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.HealthHandler)
	mux.HandleFunc("/v1/orders", handler.CreateOrderHandler)

	return middleware.RequestID(mux)
}
