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

func (c *CatalogClient) GetAllBooks(ctx context.Context, req *bookstorepb.GetAllBooksRequest) (*bookstorepb.GetAllBooksResponse, error) {
	return c.client.GetAllBooks(ctx, req)
}
