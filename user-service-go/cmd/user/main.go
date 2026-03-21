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

	bookstorepb "github.com/shank/bookstore-microservices/generated-go/bookstore"
	grpcadapter "github.com/shank/bookstore-microservices/user-service-go/internal/grpc"
	"github.com/shank/bookstore-microservices/user-service-go/internal/repository"
	"github.com/shank/bookstore-microservices/user-service-go/internal/service"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()
	port := getenv("USER_GRPC_PORT", "50051")
	dbPath := getenv("USER_DB_PATH", "user-service.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("db close error: %v", closeErr)
		}
	}()

	repo := repository.NewPostgresUserRepo(db)
	if err := repo.Init(ctx); err != nil {
		log.Fatalf("repo init failed: %v", err)
	}

	userService := service.NewUserService(repo)
	handler := grpcadapter.NewHandler(userService)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookstorepb.RegisterUserServiceServer(grpcServer, handler)

	go func() {
		log.Printf("user-service gRPC listening on :%s", port)
		if serveErr := grpcServer.Serve(listener); serveErr != nil {
			log.Fatalf("gRPC serve failed: %v", serveErr)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	gracefulStop(grpcServer)
	log.Println("user-service stopped")
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
