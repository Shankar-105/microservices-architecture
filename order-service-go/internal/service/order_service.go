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
	ErrPaymentDeclined   = errors.New("payment declined")
	ErrDependencyFailure = errors.New("dependency failure")
)

type OrderService struct {
	repo          *repository.PostgresOrderRepo
	userClient    *clients.UserClient
	catalogClient *clients.CatalogClient
	paymentClient *clients.PaymentClient
}

func NewOrderService(
	repo *repository.PostgresOrderRepo,
	userClient *clients.UserClient,
	catalogClient *clients.CatalogClient,
	paymentClient *clients.PaymentClient,
) *OrderService {
	return &OrderService{
		repo:          repo,
		userClient:    userClient,
		catalogClient: catalogClient,
		paymentClient: paymentClient,
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

	paymentResp, err := s.paymentClient.Authorize(ctx, order.OrderID, totalCents)
	if err != nil {
		if updateErr := s.repo.UpdateResult(ctx, order.OrderID, totalCents, "failed"); updateErr != nil {
			return nil, updateErr
		}
		order.Status = "failed"
		return nil, ErrDependencyFailure
	}

	if !paymentResp.GetApproved() {
		if err := s.repo.UpdateResult(ctx, order.OrderID, totalCents, "failed"); err != nil {
			return nil, err
		}
		order.Status = "failed"
		return order, ErrPaymentDeclined
	}

	if err := s.repo.UpdateResult(ctx, order.OrderID, totalCents, "paid"); err != nil {
		return nil, err
	}
	order.Status = "paid"

	return order, nil
}
