package grpc

import (
	"context"

	bookstorepb "microservices/generated-go/bookstore"
	"microservices/user-service-go/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	bookstorepb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewHandler(svc *service.UserService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) Health(_ context.Context, _ *bookstorepb.HealthCheckRequest) (*bookstorepb.HealthCheckResponse, error) {
	return &bookstorepb.HealthCheckResponse{Status: "ok", Service: "user-service-go"}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *bookstorepb.GetUserRequest) (*bookstorepb.GetUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := h.service.GetUser(ctx, req.GetUserId())
	if err != nil {
		if err == service.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &bookstorepb.GetUserResponse{
		UserId:      user.UserID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}, nil
}
