package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TrafficByPolicy represents traffic metrics for a policy
type TrafficByPolicy struct {
	Policy string `json:"policy"`
	Count  int    `json:"count"`
}

// AnalyticsData represents the analytics data
type AnalyticsData struct {
	TotalQueries    int               `json:"totalQueries"`
	BlockedRequests int               `json:"blockedRequests"`
	TopDomains      []string          `json:"topDomains"`
	TrafficByPolicy []TrafficByPolicy `json:"trafficByPolicy"`
}

// ResponseAnalytics represents the analytics response
type ResponseAnalytics struct {
	Status string        `json:"status"`
	Data   AnalyticsData `json:"data"`
	Error  *string       `json:"error,omitempty"`
}

// GetAnalyticsSummary handles GET /analytics
func (h *APIHandler) GetAnalyticsSummary(c *gin.Context) {
	// TODO: Implement logic to fetch analytics from database
	c.JSON(http.StatusOK, ResponseAnalytics{
		Status: "success",
		Data: AnalyticsData{
			TotalQueries:    0,
			BlockedRequests: 0,
			TopDomains:      []string{},
			TrafficByPolicy: []TrafficByPolicy{},
		},
	})
}
