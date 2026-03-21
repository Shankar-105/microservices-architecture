package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"microservices/gateway-go/internal/clients"
	bookstorepb "microservices/generated-go/bookstore"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	orderClient *clients.OrderClient
}

func NewHandler(orderClient *clients.OrderClient) *Handler {
	return &Handler{orderClient: orderClient}
}

func (h *Handler) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "gateway",
		"status":  "ok",
	})
}

func (h *Handler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateOrderHandler(w, r)
	case http.MethodGet:
		h.GetOrdersHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

type createOrderHTTPRequest struct {
	UserID   string `json:"user_id"`
	BookID   string `json:"book_id"`
	Quantity int32  `json:"quantity"`
}

func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var reqBody createOrderHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}

	if reqBody.UserID == "" || reqBody.BookID == "" || reqBody.Quantity <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id, book_id and positive quantity are required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.orderClient.CreateOrder(ctx, &bookstorepb.CreateOrderRequest{
		UserId:   reqBody.UserID,
		BookId:   reqBody.BookID,
		Quantity: reqBody.Quantity,
	})
	if err != nil {
		code := http.StatusInternalServerError
		message := "internal gateway error"

		if s, ok := status.FromError(err); ok {
			code = grpcCodeToHTTP(s.Code())
			message = s.Message()
		}

		writeJSON(w, code, map[string]string{"error": message})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"order_id":    grpcResp.GetOrderId(),
		"status":      grpcResp.GetStatus(),
		"total_cents": grpcResp.GetTotalCents(),
	})
}

func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.orderClient.GetOrders(ctx, &bookstorepb.GetOrdersRequest{UserId: userID})
	if err != nil {
		code := http.StatusInternalServerError
		message := "internal gateway error"

		if s, ok := status.FromError(err); ok {
			code = grpcCodeToHTTP(s.Code())
			message = s.Message()
		}

		writeJSON(w, code, map[string]string{"error": message})
		return
	}

	orders := make([]map[string]any, 0, len(grpcResp.GetOrders()))
	for _, order := range grpcResp.GetOrders() {
		orders = append(orders, map[string]any{
			"order_id":    order.GetOrderId(),
			"user_id":     order.GetUserId(),
			"book_id":     order.GetBookId(),
			"quantity":    order.GetQuantity(),
			"total_cents": order.GetTotalCents(),
			"status":      order.GetStatus(),
			"created_at":  order.GetCreatedAt(),
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"orders": orders})
}

func grpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.FailedPrecondition:
		return http.StatusConflict
	case codes.Unavailable, codes.DeadlineExceeded:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
