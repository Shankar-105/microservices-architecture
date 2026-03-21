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
	mux.HandleFunc("/createorder", handler.CreateOrderHandler)
	mux.HandleFunc("/getorders", handler.GetOrdersHandler)

	return middleware.RequestID(mux)
}
