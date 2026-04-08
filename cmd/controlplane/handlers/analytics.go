package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type QueryLogEntry struct {
	ID              uint      `json:"id"`
	Domain          string    `json:"domain"`
	ClientIP        string    `json:"client_ip"`
	Action          string    `json:"action"`
	Timestamp       time.Time `json:"timestamp"`
	IsSuspicious    bool      `json:"is_suspicious"`
	ThreatScore     float64   `json:"threat_score"`
	DetectionMethod string    `json:"detection_method,omitempty"`
	ThreatReason    string    `json:"threat_reason,omitempty"`
}

type AnalyticsSummaryData struct {
	TotalQueries     uint64  `json:"total_queries"`
	BlockedQueries   uint64  `json:"blocked_queries"`
	AllowedQueries   uint64  `json:"allowed_queries"`
	BlockRatePercent float64 `json:"block_rate_percent"`
}

type ResponseAnalyticsSummary struct {
	Status string               `json:"status"`
	Data   AnalyticsSummaryData `json:"data"`
	Error  *string              `json:"error"`
}

type ResponseQueryLogList struct {
	Status string          `json:"status"`
	Data   []QueryLogEntry `json:"data"`
	Error  *string         `json:"error"`
}

// GetAnalyticsSummary handles GET /analytics/summary
func (h *APIHandler) GetAnalyticsSummary(c *gin.Context) {
	stats, err := h.Store.Statistics.ListRecent(1)
	if err != nil || len(stats) == 0 {
		c.JSON(http.StatusOK, ResponseAnalyticsSummary{
			Status: "success",
			Data:   AnalyticsSummaryData{},
		})
		return
	}

	s := stats[0]
	var blockRate float64
	if s.TotalQueries > 0 {
		blockRate = float64(s.BlockedQueries) / float64(s.TotalQueries) * 100
	}

	c.JSON(http.StatusOK, ResponseAnalyticsSummary{
		Status: "success",
		Data: AnalyticsSummaryData{
			TotalQueries:     s.TotalQueries,
			BlockedQueries:   s.BlockedQueries,
			AllowedQueries:   s.AllowedQueries,
			BlockRatePercent: blockRate,
		},
	})
}

// GetAuditLogs handles GET /analytics/audits
// Returns recent DNS query logs as the audit trail.
func (h *APIHandler) GetAuditLogs(c *gin.Context) {
	queries, err := h.Store.QueryLogs.ListRecent(100)
	if err != nil {
		errMsg := "failed to fetch query logs"
		c.JSON(http.StatusInternalServerError, ResponseQueryLogList{Status: "error", Error: &errMsg})
		return
	}

	entries := make([]QueryLogEntry, 0, len(queries))
	for _, q := range queries {
		entries = append(entries, QueryLogEntry{
			ID:              q.ID,
			Domain:          q.Domain,
			ClientIP:        q.ClientIP,
			Action:          q.Action,
			Timestamp:       q.Timestamp,
			IsSuspicious:    q.IsSuspicious,
			ThreatScore:     q.ThreatScore,
			DetectionMethod: q.DetectionMethod,
			ThreatReason:    q.ThreatReason,
		})
	}

	c.JSON(http.StatusOK, ResponseQueryLogList{
		Status: "success",
		Data:   entries,
	})
}
