package authclient

import (
	"google.golang.org/grpc"

	"github.com/vasapolrittideah/optimize-api/shared/discovery"
	authpbv1 "github.com/vasapolrittideah/optimize-api/shared/protos/auth/v1"
)

type AuthServiceClient struct {
	Client authpbv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient(serviceName string, consulRegistry *discovery.ConsulRegistry) (*AuthServiceClient, error) {
	conn, err := consulRegistry.Connect(serviceName)
	if err != nil {
		return nil, err
	}

	client := authpbv1.NewAuthServiceClient(conn)

	return &AuthServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *AuthServiceClient) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
