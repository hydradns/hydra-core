package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

// Policy represents a policy in API responses
type Policy struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Action      string   `json:"action"`
	RedirectIP  string   `json:"redirect_ip,omitempty"`
	Domains     []string `json:"domains"`
	Priority    int      `json:"priority"`
	Enabled     bool     `json:"enabled"`
}

type PolicyListData struct {
	TotalPolicies    int      `json:"total_policies"`
	ActivePolicies   int      `json:"active_policies"`
	InactivePolicies int      `json:"inactive_policies"`
	List             []Policy `json:"list"`
}

type ResponsePolicyList struct {
	Status string         `json:"status"`
	Data   PolicyListData `json:"data"`
	Error  *string        `json:"error"`
}

type ResponsePolicySingle struct {
	Status string  `json:"status"`
	Data   Policy  `json:"data"`
	Error  *string `json:"error"`
}

type CreatePolicyRequest struct {
	ID          string   `json:"id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Action      string   `json:"action" binding:"required"`
	RedirectIP  string   `json:"redirect_ip"`
	Domains     []string `json:"domains" binding:"required"`
	Priority    int      `json:"priority"`
}

func policyFromModel(m models.Policy) Policy {
	var domains []string
	if m.Domains != "" {
		_ = json.Unmarshal([]byte(m.Domains), &domains)
	}
	return Policy{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Action:      m.Action,
		RedirectIP:  m.RedirectIP,
		Domains:     domains,
		Priority:    m.Priority,
		Enabled:     m.Enabled,
	}
}

// ListPolicies handles GET /policies
func (h *APIHandler) ListPolicies(c *gin.Context) {
	models, err := h.Store.Policies.List()
	if err != nil {
		errMsg := "failed to fetch policies"
		c.JSON(http.StatusInternalServerError, ResponsePolicyList{Status: "error", Error: &errMsg})
		return
	}

	var list []Policy
	activeCount := 0
	for _, m := range models {
		p := policyFromModel(m)
		list = append(list, p)
		if p.Enabled {
			activeCount++
		}
	}

	c.JSON(http.StatusOK, ResponsePolicyList{
		Status: "success",
		Data: PolicyListData{
			TotalPolicies:    len(list),
			ActivePolicies:   activeCount,
			InactivePolicies: len(list) - activeCount,
			List:             list,
		},
	})
}

// GetPolicy handles GET /policies/:id
func (h *APIHandler) GetPolicy(c *gin.Context) {
	m, err := h.Store.Policies.GetByID(c.Param("id"))
	if err != nil {
		errMsg := "policy not found"
		c.JSON(http.StatusNotFound, ResponsePolicySingle{Status: "error", Error: &errMsg})
		return
	}
	c.JSON(http.StatusOK, ResponsePolicySingle{
		Status: "success",
		Data:   policyFromModel(*m),
	})
}

// CreatePolicy handles POST /policies
func (h *APIHandler) CreatePolicy(c *gin.Context) {
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, ResponsePolicySingle{Status: "error", Error: &errMsg})
		return
	}

	domainsJSON, _ := json.Marshal(req.Domains)

	m := &models.Policy{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Action:      req.Action,
		RedirectIP:  req.RedirectIP,
		Domains:     string(domainsJSON),
		Priority:    req.Priority,
		Enabled:     true,
	}
	if err := h.Store.Policies.Create(m); err != nil {
		errMsg := "failed to create policy"
		c.JSON(http.StatusInternalServerError, ResponsePolicySingle{Status: "error", Error: &errMsg})
		return
	}
	c.JSON(http.StatusCreated, ResponsePolicySingle{
		Status: "success",
		Data:   policyFromModel(*m),
	})
}

// DeletePolicy handles DELETE /policies/:id
func (h *APIHandler) DeletePolicy(c *gin.Context) {
	if err := h.Store.Policies.Delete(c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		errMsg := "failed to delete policy"
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
			errMsg = "policy not found"
		}
		c.JSON(status, ResponseGeneric{Status: "error", Error: &errMsg})
		return
	}
	c.JSON(http.StatusOK, ResponseGeneric{Status: "success", Data: map[string]interface{}{}})
}
