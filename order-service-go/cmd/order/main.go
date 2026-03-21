package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	bookstorepb "microservices/generated-go/bookstore"
	"microservices/order-service-go/internal/clients"
	grpcadapter "microservices/order-service-go/internal/grpc"
	"microservices/order-service-go/internal/repository"
	"microservices/order-service-go/internal/service"

	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()
	port := getenv("ORDER_GRPC_PORT", "50052")
	dbPath := getenv("ORDER_DB_PATH", "order-service.db")
	userServiceAddr := getenv("USER_SERVICE_ADDR", "localhost:50051")
	catalogServiceAddr := getenv("CATALOG_SERVICE_ADDR", "localhost:50053")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("db close error: %v", closeErr)
		}
	}()

	repo := repository.NewPostgresOrderRepo(db)
	if err := repo.Init(ctx); err != nil {
		log.Fatalf("repo init failed: %v", err)
	}

	userClient, err := clients.NewUserClient(ctx, userServiceAddr)
	if err != nil {
		log.Fatalf("user client init failed: %v", err)
	}
	defer func() {
		if closeErr := userClient.Close(); closeErr != nil {
			log.Printf("user client close error: %v", closeErr)
		}
	}()

	catalogClient, err := clients.NewCatalogClient(ctx, catalogServiceAddr)
	if err != nil {
		log.Fatalf("catalog client init failed: %v", err)
	}
	defer func() {
		if closeErr := catalogClient.Close(); closeErr != nil {
			log.Printf("catalog client close error: %v", closeErr)
		}
	}()

	orderService := service.NewOrderService(repo, userClient, catalogClient)
	handler := grpcadapter.NewHandler(orderService)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookstorepb.RegisterOrderServiceServer(grpcServer, handler)

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
