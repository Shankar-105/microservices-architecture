package grpc

import (
	"context"
	"testing"

	bookstorepb "microservices/generated-go/bookstore"
	"microservices/order-service-go/internal/service"

	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateOrderRejectsNilRequest(t *testing.T) {
	handler := NewHandler(&service.OrderService{})
	_, err := handler.CreateOrder(context.Background(), nil)
	if status.Code(err) != gcodes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestCreateOrderRejectsQuantityBoundary(t *testing.T) {
	handler := NewHandler(&service.OrderService{})
	for _, quantity := range []int32{0, -1, 10001} {
		request := &bookstorepb.CreateOrderRequest{
			UserId: "u-100", BookId: "b-100", Quantity: quantity,
		}
		_, err := handler.CreateOrder(context.Background(), request)
		if status.Code(err) != gcodes.InvalidArgument {
			t.Fatalf("quantity %d: expected InvalidArgument, got %v", quantity, status.Code(err))
		}
	}
}

func TestGetOrdersRejectsOversizedUserID(t *testing.T) {
	handler := NewHandler(&service.OrderService{})
	request := &bookstorepb.GetOrdersRequest{UserId: string(make([]byte, 129))}
	_, err := handler.GetOrders(context.Background(), request)
	if status.Code(err) != gcodes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}
