package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Blocklist represents a blocklist entry
type Blocklist struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Category          string `json:"category"`
	DomainsCount      int    `json:"domains_count"`
	AvgBlockedQueries int    `json:"avg_blocked_queries"`
	Active            bool   `json:"active"`
}

// BlocklistListData represents blocklist list data with statistics
type BlocklistListData struct {
	TotalBlocklists   int         `json:"total_blocklists"`
	TotalDomains      int         `json:"total_domains"`
	AvgBlockedQueries float64     `json:"avg_blocked_queries"`
	ActiveLists       []Blocklist `json:"active_lists"`
}

// ResponseBlocklistList represents a list of blocklists response
type ResponseBlocklistList struct {
	Status string            `json:"status"`
	Data   BlocklistListData `json:"data"`
	Error  *string           `json:"error"`
}

// ResponseBlocklistSingle represents a single blocklist entry response
type ResponseBlocklistSingle struct {
	Status string    `json:"status"`
	Data   Blocklist `json:"data"`
	Error  *string   `json:"error"`
}

var mockBlocklists = []Blocklist{
	{ID: "1", Name: "Ads Block", Category: "ads", DomainsCount: 10000, AvgBlockedQueries: 1500, Active: true},
	{ID: "2", Name: "Malware Block", Category: "malware", DomainsCount: 5000, AvgBlockedQueries: 1200, Active: true},
	{ID: "3", Name: "Tracking Block", Category: "tracking", DomainsCount: 3000, AvgBlockedQueries: 500, Active: true},
}

// ListBlocklists handles GET /blocklists
func (h *APIHandler) ListBlocklists(c *gin.Context) {
	// TODO: Implement logic to fetch blocklists from database
	totalDomains := 0
	totalBlocked := 0
	for _, bl := range mockBlocklists {
		totalDomains += bl.DomainsCount
		totalBlocked += bl.AvgBlockedQueries
	}
	avgBlocked := 0.0
	if len(mockBlocklists) > 0 {
		avgBlocked = float64(totalBlocked) / float64(len(mockBlocklists))
	}

	c.JSON(http.StatusOK, ResponseBlocklistList{
		Status: "success",
		Data: BlocklistListData{
			TotalBlocklists:   len(mockBlocklists),
			TotalDomains:      totalDomains,
			AvgBlockedQueries: avgBlocked,
			ActiveLists:       mockBlocklists,
		},
		Error: nil,
	})
}

// GetBlocklist handles GET /blocklists/:id
func (h *APIHandler) GetBlocklist(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to fetch blocklist from database
	for _, blocklist := range mockBlocklists {
		if blocklist.ID == id {
			c.JSON(http.StatusOK, ResponseBlocklistSingle{
				Status: "success",
				Data:   blocklist,
				Error:  nil,
			})
			return
		}
	}
	errMsg := "blocklist not found"
	c.JSON(http.StatusNotFound, ResponseBlocklistSingle{
		Status: "error",
		Data:   Blocklist{},
		Error:  &errMsg,
	})
}

// CreateBlocklist handles POST /blocklists
func (h *APIHandler) CreateBlocklist(c *gin.Context) {
	var blocklist Blocklist
	if err := c.ShouldBindJSON(&blocklist); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, ResponseBlocklistSingle{
			Status: "error",
			Data:   Blocklist{},
			Error:  &errMsg,
		})
		return
	}
	// TODO: Implement logic to add blocklist entry to database
	mockBlocklists = append(mockBlocklists, blocklist)
	c.JSON(http.StatusCreated, ResponseBlocklistSingle{
		Status: "success",
		Data:   blocklist,
		Error:  nil,
	})
}

// DeleteBlocklist handles DELETE /blocklists/:id
func (h *APIHandler) DeleteBlocklist(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to delete blocklist from database
	for i, blocklist := range mockBlocklists {
		if blocklist.ID == id {
			mockBlocklists = append(mockBlocklists[:i], mockBlocklists[i+1:]...)
			c.JSON(http.StatusOK, ResponseGeneric{
				Status: "success",
				Data:   map[string]interface{}{},
				Error:  nil,
			})
			return
		}
	}
	errMsg := "blocklist not found"
	c.JSON(http.StatusNotFound, ResponseGeneric{
		Status: "error",
		Data:   nil,
		Error:  &errMsg,
	})
}
