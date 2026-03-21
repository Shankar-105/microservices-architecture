package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"microservices/gateway-go/internal/clients"
	httptransport "microservices/gateway-go/internal/http"
)

func main() {
	port := getenv("GATEWAY_HTTP_PORT", "8080")
	orderServiceAddr := getenv("ORDER_SERVICE_ADDR", "localhost:50052")

	orderClient, err := clients.NewOrderClient(context.Background(), orderServiceAddr)
	if err != nil {
		log.Fatalf("order client init failed: %v", err)
	}
	defer func() {
		if closeErr := orderClient.Close(); closeErr != nil {
			log.Printf("order client close error: %v", closeErr)
		}
	}()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httptransport.NewRouter(orderClient),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("gateway listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gateway failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("gateway shutdown error: %v", err)
	}
	log.Println("gateway stopped")
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
