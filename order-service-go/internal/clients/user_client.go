package clients

import (
	"context"

	bookstorepb "microservices/generated-go/bookstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client bookstorepb.UserServiceClient
}

func NewUserClient(ctx context.Context, address string) (*UserClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		conn:   conn,
		client: bookstorepb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*bookstorepb.GetUserResponse, error) {
	return c.client.GetUser(ctx, &bookstorepb.GetUserRequest{UserId: userID})
}
