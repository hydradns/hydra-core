package handlers

import "github.com/lopster568/phantomDNS/internal/grpc/client"

type APIHandler struct {
	GRPCClient *client.Client
}

func NewAPIHandler(grpcClient *client.Client) *APIHandler {
	return &APIHandler{
		GRPCClient: grpcClient,
	}
}
