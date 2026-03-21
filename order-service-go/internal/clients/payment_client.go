package clients

import (
	"context"

	bookstorepb "microservices/generated-go/bookstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	conn   *grpc.ClientConn
	client bookstorepb.PaymentServiceClient
}

func NewPaymentClient(ctx context.Context, address string) (*PaymentClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &PaymentClient{
		conn:   conn,
		client: bookstorepb.NewPaymentServiceClient(conn),
	}, nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

func (c *PaymentClient) Authorize(ctx context.Context, orderID string, amountCents int64) (*bookstorepb.AuthorizePaymentResponse, error) {
	return c.client.Authorize(ctx, &bookstorepb.AuthorizePaymentRequest{OrderId: orderID, AmountCents: amountCents})
}
