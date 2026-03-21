package clients

import (
	"context"

	bookstorepb "microservices/generated-go/bookstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	conn   *grpc.ClientConn
	client bookstorepb.OrderServiceClient
}

func NewOrderClient(ctx context.Context, address string) (*OrderClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &OrderClient{
		conn:   conn,
		client: bookstorepb.NewOrderServiceClient(conn),
	}, nil
}

func (c *OrderClient) Close() error {
	return c.conn.Close()
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *bookstorepb.CreateOrderRequest) (*bookstorepb.CreateOrderResponse, error) {
	return c.client.CreateOrder(ctx, req)
}
