package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lopster568/phantomDNS/internal/storage/models"
)

// Blocklist represents a blocklist source in API responses
type Blocklist struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Format       string    `json:"format"`
	Category     string    `json:"category"`
	DomainsCount int64     `json:"domains_count"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BlocklistListData struct {
	TotalBlocklists int         `json:"total_blocklists"`
	TotalDomains    int64       `json:"total_domains"`
	ActiveLists     []Blocklist `json:"active_lists"`
}

type ResponseBlocklistList struct {
	Status string            `json:"status"`
	Data   BlocklistListData `json:"data"`
	Error  *string           `json:"error"`
}

type ResponseBlocklistSingle struct {
	Status string    `json:"status"`
	Data   Blocklist `json:"data"`
	Error  *string   `json:"error"`
}

type CreateBlocklistRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Format   string `json:"format" binding:"required"`
	Category string `json:"category"`
}

func blocklistFromSource(src models.BlocklistSource, count int64) Blocklist {
	return Blocklist{
		ID:           src.ID,
		Name:         src.Name,
		URL:          src.URL,
		Format:       src.Format,
		Category:     src.Category,
		DomainsCount: count,
		Enabled:      src.Enabled,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

// ListBlocklists handles GET /blocklists
func (h *APIHandler) ListBlocklists(c *gin.Context) {
	sources, err := h.Store.Blocklist.ListSources()
	if err != nil {
		errMsg := "failed to fetch blocklist sources"
		c.JSON(http.StatusInternalServerError, ResponseBlocklistList{Status: "error", Error: &errMsg})
		return
	}

	counts, err := h.Store.Blocklist.CountEntriesGroupedBySource()
	if err != nil {
		counts = map[string]int64{}
	}

	var lists []Blocklist
	var totalDomains int64
	for _, src := range sources {
		count := counts[src.ID]
		totalDomains += count
		lists = append(lists, blocklistFromSource(src, count))
	}

	c.JSON(http.StatusOK, ResponseBlocklistList{
		Status: "success",
		Data: BlocklistListData{
			TotalBlocklists: len(lists),
			TotalDomains:    totalDomains,
			ActiveLists:     lists,
		},
	})
}

// GetBlocklist handles GET /blocklists/:id
func (h *APIHandler) GetBlocklist(c *gin.Context) {
	src, err := h.Store.Blocklist.GetSource(c.Param("id"))
	if err != nil {
		errMsg := "blocklist not found"
		c.JSON(http.StatusNotFound, ResponseBlocklistSingle{Status: "error", Error: &errMsg})
		return
	}
	count, _ := h.Store.Blocklist.CountEntriesBySource(src.ID)
	c.JSON(http.StatusOK, ResponseBlocklistSingle{
		Status: "success",
		Data:   blocklistFromSource(*src, count),
	})
}

// CreateBlocklist handles POST /blocklists
func (h *APIHandler) CreateBlocklist(c *gin.Context) {
	var req CreateBlocklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, ResponseBlocklistSingle{Status: "error", Error: &errMsg})
		return
	}

	src := &models.BlocklistSource{
		ID:        req.ID,
		Name:      req.Name,
		URL:       req.URL,
		Format:    req.Format,
		Category:  req.Category,
		Enabled:   true,
		CreatedAt: time.Now(),
	}
	if err := h.Store.Blocklist.CreateSource(src); err != nil {
		errMsg := "failed to create blocklist source"
		c.JSON(http.StatusInternalServerError, ResponseBlocklistSingle{Status: "error", Error: &errMsg})
		return
	}

	c.JSON(http.StatusCreated, ResponseBlocklistSingle{
		Status: "success",
		Data:   blocklistFromSource(*src, 0),
	})
}

// DeleteBlocklist handles DELETE /blocklists/:id
func (h *APIHandler) DeleteBlocklist(c *gin.Context) {
	if err := h.Store.Blocklist.DeleteSource(c.Param("id")); err != nil {
		errMsg := "blocklist not found"
		c.JSON(http.StatusNotFound, ResponseGeneric{Status: "error", Error: &errMsg})
		return
	}
	c.JSON(http.StatusOK, ResponseGeneric{Status: "success", Data: map[string]interface{}{}})
}
