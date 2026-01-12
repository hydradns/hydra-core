package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DnsEngineStatusData represents DNS engine status
type DnsEngineStatusData struct {
	Enabled        bool    `json:"enabled"`
	AvgQueryTimeMs float64 `json:"avg_query_time_ms"`
	QueryRateQps   float64 `json:"query_rate_qps"`
}

// ResponseDnsEngineStatus represents DNS engine status response
type ResponseDnsEngineStatus struct {
	Status string              `json:"status"`
	Data   DnsEngineStatusData `json:"data"`
	Error  *string             `json:"error"`
}

// Resolver represents an upstream DNS resolver
type Resolver struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
}

// ResponseResolverList represents a list of resolvers response
type ResponseResolverList struct {
	Status string     `json:"status"`
	Data   []Resolver `json:"data"`
	Error  *string    `json:"error"`
}

// ResponseResolverSingle represents a single resolver response
type ResponseResolverSingle struct {
	Status string   `json:"status"`
	Data   Resolver `json:"data"`
	Error  *string  `json:"error"`
}

// ToggleDnsEngineRequest represents request to toggle DNS engine
type ToggleDnsEngineRequest struct {
	Enabled bool `json:"enabled"`
}

var mockResolvers = []Resolver{
	{ID: "1", Name: "Google DNS", Address: "8.8.8.8", Protocol: "udp"},
	{ID: "2", Name: "Cloudflare DNS", Address: "1.1.1.1", Protocol: "udp"},
}

var dnsEngineEnabled = true

// GetDnsEngineStatus handles GET /dns/engine
func (h *APIHandler) GetDnsEngineStatus(c *gin.Context) {
	// TODO: Implement logic to fetch DNS engine status from dataplane
	c.JSON(http.StatusOK, ResponseDnsEngineStatus{
		Status: "success",
		Data: DnsEngineStatusData{
			Enabled:        dnsEngineEnabled,
			AvgQueryTimeMs: 10.4,
			QueryRateQps:   350.0,
		},
		Error: nil,
	})
}

// ToggleDnsEngine handles POST /dns/engine
func (h *APIHandler) ToggleDnsEngine(c *gin.Context) {
	var req ToggleDnsEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, ResponseGeneric{
			Status: "error",
			Error:  &errMsg,
		})
		return
	}

	// 1. Persist desired state (source of truth)
	if err := h.Store.SystemState.SetDNSEnabled(req.Enabled); err != nil {
		errMsg := "failed to persist DNS engine state"
		c.JSON(http.StatusInternalServerError, ResponseGeneric{
			Status: "error",
			Error:  &errMsg,
		})
		return
	}

	// 2. Apply desired state to dataplane via gRPC
	if err := h.DataPlaneClient.SetAcceptQueries(req.Enabled); err != nil {
		errMsg := "failed to apply DNS engine state to dataplane"
		c.JSON(http.StatusBadGateway, ResponseGeneric{
			Status: "error",
			Error:  &errMsg,
		})
		return
	}

	// 3. Respond with acknowledged intent
	c.JSON(http.StatusOK, ResponseGeneric{
		Status: "success",
		Data: map[string]interface{}{
			"enabled": req.Enabled,
		},
	})
}

// ListResolvers handles GET /dns/resolvers
func (h *APIHandler) ListResolvers(c *gin.Context) {
	// TODO: Implement logic to fetch resolvers from dataplane
	c.JSON(http.StatusOK, ResponseResolverList{
		Status: "success",
		Data:   mockResolvers,
		Error:  nil,
	})
}

// AddResolver handles POST /dns/resolvers
func (h *APIHandler) AddResolver(c *gin.Context) {
	var resolver Resolver
	if err := c.ShouldBindJSON(&resolver); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, ResponseResolverSingle{
			Status: "error",
			Data:   Resolver{},
			Error:  &errMsg,
		})
		return
	}
	// TODO: Implement logic to add resolver via gRPC
	mockResolvers = append(mockResolvers, resolver)
	c.JSON(http.StatusCreated, ResponseResolverSingle{
		Status: "success",
		Data:   resolver,
		Error:  nil,
	})
}

// DeleteResolver handles DELETE /dns/resolvers/:id
func (h *APIHandler) DeleteResolver(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to delete resolver via gRPC
	for i, resolver := range mockResolvers {
		if resolver.ID == id {
			mockResolvers = append(mockResolvers[:i], mockResolvers[i+1:]...)
			c.JSON(http.StatusOK, ResponseGeneric{
				Status: "success",
				Data:   map[string]interface{}{},
				Error:  nil,
			})
			return
		}
	}
	errMsg := "resolver not found"
	c.JSON(http.StatusNotFound, ResponseGeneric{
		Status: "error",
		Data:   nil,
		Error:  &errMsg,
	})
}
