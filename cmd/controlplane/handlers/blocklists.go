package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Blocklist represents a blocklist entry
type Blocklist struct {
	ID       string `json:"id"`
	Domain   string `json:"domain"`
	Category string `json:"category"`
}

// ResponseBlocklistList represents a list of blocklists response
type ResponseBlocklistList struct {
	Status string      `json:"status"`
	Data   []Blocklist `json:"data"`
	Error  *string     `json:"error,omitempty"`
}

// ResponseBlocklistSingle represents a single blocklist entry response
type ResponseBlocklistSingle struct {
	Status string    `json:"status"`
	Data   Blocklist `json:"data"`
	Error  *string   `json:"error,omitempty"`
}

// ListBlocklists handles GET /blocklists
func (h *APIHandler) ListBlocklists(c *gin.Context) {
	// TODO: Implement logic to fetch blocklists from database
	c.JSON(http.StatusOK, ResponseBlocklistList{
		Status: "success",
		Data:   []Blocklist{},
	})
}

// AddToBlocklist handles POST /blocklists
func (h *APIHandler) AddToBlocklist(c *gin.Context) {
	var blocklist Blocklist
	if err := c.ShouldBindJSON(&blocklist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Implement logic to add blocklist entry to database
	c.JSON(http.StatusCreated, ResponseBlocklistSingle{
		Status: "success",
		Data:   blocklist,
	})
}
