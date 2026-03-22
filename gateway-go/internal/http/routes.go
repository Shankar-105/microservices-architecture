package http

import (
	"net/http"

	"microservices/gateway-go/internal/clients"
	"microservices/gateway-go/internal/middleware"
)

func NewRouter(orderClient *clients.OrderClient, catalogClient *clients.CatalogClient) http.Handler {
	handler := NewHandler(orderClient, catalogClient)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.HealthHandler)
	mux.HandleFunc("/createorder", handler.CreateOrderHandler)
	mux.HandleFunc("/getorders", handler.GetOrdersHandler)
	mux.HandleFunc("/getallbooks", handler.GetAllBooksHandler)

	return middleware.RequestID(mux)
}
