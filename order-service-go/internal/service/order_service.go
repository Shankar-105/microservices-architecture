package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shank/bookstore-microservices/order-service-go/internal/repository"
)

type OrderService struct {
	repo *repository.PostgresOrderRepo
}

func NewOrderService(repo *repository.PostgresOrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID, bookID string, quantity int32) (*repository.Order, error) {
	order := &repository.Order{
		OrderID:    fmt.Sprintf("ord-%d", time.Now().UnixNano()),
		UserID:     userID,
		BookID:     bookID,
		Quantity:   quantity,
		TotalCents: 0,
		Status:     "pending",
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
