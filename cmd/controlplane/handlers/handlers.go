package handlers

import (
	client "github.com/lopster568/phantomDNS/internal/grpc/controlplane"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"
)

// APIHandler contains dependencies for API endpoints
type APIHandler struct {
	Store           repositories.Store
	DataPlaneClient *client.Client
}

func NewAPIHandler(
	store repositories.Store,
	dataPlaneClient *client.Client,
) *APIHandler {
	return &APIHandler{
		Store:           store,
		DataPlaneClient: dataPlaneClient,
	}
}
