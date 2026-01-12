package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/lopster568/phantomDNS/internal/gen/proto/phantomdns/v1"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	status pb.DataPlaneStatusServiceClient
}

// New creates and connects the gRPC client.
func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // TLS later
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:   conn,
		status: pb.NewDataPlaneStatusServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// GetStatus fetches dataplane runtime status.
func (c *Client) GetStatus() (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.status.GetStatus(ctx, &pb.StatusRequest{})
}

func (c *Client) SetAcceptQueries(enabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.status.SetAcceptQueries(ctx, &pb.SetAcceptQueriesRequest{
		Enabled: enabled,
	})
	return err
}
