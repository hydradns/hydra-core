package server

import (
	"context"

	"github.com/lopster568/phantomDNS/proto/healthpb"
)

// HealthService implements the HealthServer interface.
type HealthService struct {
	healthpb.UnimplementedHealthServer
}

// NewHealthService creates a new health service.
func NewHealthService() *HealthService {
	return &HealthService{}
}

// Check implements the health check RPC.
func (h *HealthService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: "OK",
	}, nil
}
