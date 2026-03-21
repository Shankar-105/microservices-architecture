package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	port := getenv("ORDER_GRPC_PORT", "50052")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	grpcServer := grpc.NewServer()

	go func() {
		log.Printf("order-service gRPC listening on :%s", port)
		if serveErr := grpcServer.Serve(listener); serveErr != nil {
			log.Fatalf("gRPC serve failed: %v", serveErr)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	gracefulStop(grpcServer)
	log.Println("order-service stopped")
}

func gracefulStop(server *grpc.Server) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		server.Stop()
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
