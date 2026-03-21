package grpc

import (
	"context"

	bookstorepb "microservices/generated-go/bookstore"
	"microservices/order-service-go/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	bookstorepb.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewHandler(svc *service.OrderService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) Health(_ context.Context, _ *bookstorepb.HealthCheckRequest) (*bookstorepb.HealthCheckResponse, error) {
	return &bookstorepb.HealthCheckResponse{Status: "ok", Service: "order-service-go"}, nil
}

func (h *Handler) CreateOrder(ctx context.Context, req *bookstorepb.CreateOrderRequest) (*bookstorepb.CreateOrderResponse, error) {
	if req.GetUserId() == "" || req.GetBookId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and book_id are required")
	}

	if req.GetQuantity() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	order, err := h.service.CreateOrder(ctx, req.GetUserId(), req.GetBookId(), req.GetQuantity())
	if err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrBookNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case service.ErrBookUnavailable:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case service.ErrDependencyFailure:
			return nil, status.Error(codes.Unavailable, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &bookstorepb.CreateOrderResponse{
		OrderId:    order.OrderID,
		Status:     order.Status,
		TotalCents: order.TotalCents,
	}, nil
}
