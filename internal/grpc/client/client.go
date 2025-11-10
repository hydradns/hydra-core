package client

import (
	"context"
	"fmt"
	"time"

	"github.com/lopster568/phantomDNS/proto/healthpb"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	health healthpb.HealthClient
}

// New creates and connects the gRPC client.
func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // TLS later
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:   conn,
		health: healthpb.NewHealthClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// CheckHealth calls the health check RPC.
func (c *Client) CheckHealth() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := c.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
