package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceHealth represents the health status of a service
type ServiceHealth struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
}

// SystemStatusData represents the system status data
type SystemStatusData struct {
	Version  string          `json:"version"`
	Uptime   string          `json:"uptime"`
	Services []ServiceHealth `json:"services"`
}

// ResponseSystemStatus represents the system status response
type ResponseSystemStatus struct {
	Status string           `json:"status"`
	Data   SystemStatusData `json:"data"`
	Error  *string          `json:"error,omitempty"`
}

// GetSystemStatus handles GET /system/status
func (h *APIHandler) GetSystemStatus(c *gin.Context) {
	// TODO: Implement logic to fetch system status from dataplane
	c.JSON(http.StatusOK, ResponseSystemStatus{
		Status: "success",
		Data: SystemStatusData{
			Version:  "1.0.0",
			Uptime:   "0h 0m",
			Services: []ServiceHealth{},
		},
	})
}
