package clients

import (
	"context"

	bookstorepb "microservices/generated-go/bookstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CatalogClient struct {
	conn   *grpc.ClientConn
	client bookstorepb.CatalogServiceClient
}

func NewCatalogClient(ctx context.Context, address string) (*CatalogClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &CatalogClient{
		conn:   conn,
		client: bookstorepb.NewCatalogServiceClient(conn),
	}, nil
}

func (c *CatalogClient) Close() error {
	return c.conn.Close()
}

func (c *CatalogClient) GetBook(ctx context.Context, bookID string) (*bookstorepb.GetBookResponse, error) {
	return c.client.GetBook(ctx, &bookstorepb.GetBookRequest{BookId: bookID})
}
