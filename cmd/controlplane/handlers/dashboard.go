package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// QueryVolumeChart represents query volume chart data
type QueryVolumeChart struct {
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

// DashboardSummaryData represents dashboard summary metrics
type DashboardSummaryData struct {
	TotalQueries      int              `json:"total_queries"`
	BlockedQueries    int              `json:"blocked_queries"`
	AvgQueryTimeMs    float64          `json:"avg_query_time_ms"`
	BlockRatePercent  float64          `json:"block_rate_percent"`
	QueryVolumeChart  QueryVolumeChart `json:"query_volume_chart"`
	BlockedByCategory map[string]int   `json:"blocked_by_category"`
	TopBlockedDomains []string         `json:"top_blocked_domains"`
	TopQueriedDomains []string         `json:"top_queried_domains"`
}

// ResponseDashboardSummary represents dashboard summary response
type ResponseDashboardSummary struct {
	Status string               `json:"status"`
	Data   DashboardSummaryData `json:"data"`
	Error  *string              `json:"error"`
}

// GetDashboardSummary handles GET /dashboard/summary
func (h *APIHandler) GetDashboardSummary(c *gin.Context) {
	// TODO: Implement logic to fetch dashboard metrics from database
	c.JSON(http.StatusOK, ResponseDashboardSummary{
		Status: "success",
		Data: DashboardSummaryData{
			TotalQueries:     125000,
			BlockedQueries:   3200,
			AvgQueryTimeMs:   12.3,
			BlockRatePercent: 2.56,
			QueryVolumeChart: QueryVolumeChart{
				Labels: []string{"00h", "06h", "12h", "18h"},
				Values: []int{15000, 32000, 48000, 30000},
			},
			BlockedByCategory: map[string]int{
				"malware":  1200,
				"ads":      1500,
				"tracking": 500,
			},
			TopBlockedDomains: []string{"ad.doubleclick.net", "malicious.com", "tracking.net"},
			TopQueriedDomains: []string{"google.com", "cloudflare.com", "netflix.com"},
		},
		Error: nil,
	})
}
