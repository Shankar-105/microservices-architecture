package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"microservices/order-service-go/internal/clients"
	"microservices/order-service-go/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrBookNotFound      = errors.New("book not found")
	ErrBookUnavailable   = errors.New("book unavailable")
	ErrDependencyFailure = errors.New("dependency failure")
)

type OrderService struct {
	repo          *repository.PostgresOrderRepo
	userClient    *clients.UserClient
	catalogClient *clients.CatalogClient
}

func NewOrderService(
	repo *repository.PostgresOrderRepo,
	userClient *clients.UserClient,
	catalogClient *clients.CatalogClient,
) *OrderService {
	return &OrderService{
		repo:          repo,
		userClient:    userClient,
		catalogClient: catalogClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID, bookID string, quantity int32) (*repository.Order, error) {
	if _, err := s.userClient.GetUser(ctx, userID); err != nil {
		if grpcStatus, ok := status.FromError(err); ok && grpcStatus.Code() == codes.NotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrDependencyFailure
	}

	bookResp, err := s.catalogClient.GetBook(ctx, bookID)
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok && grpcStatus.Code() == codes.NotFound {
			return nil, ErrBookNotFound
		}
		return nil, ErrDependencyFailure
	}

	if !bookResp.GetAvailable() {
		return nil, ErrBookUnavailable
	}

	totalCents := bookResp.GetPriceCents() * int64(quantity)

	order := &repository.Order{
		OrderID:    fmt.Sprintf("ord-%d", time.Now().UnixNano()),
		UserID:     userID,
		BookID:     bookID,
		Quantity:   quantity,
		TotalCents: totalCents,
		Status:     "pending",
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateResult(ctx, order.OrderID, totalCents, "confirmed"); err != nil {
		return nil, err
	}
	order.Status = "confirmed"

	return order, nil
}
