package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Policy represents a DNS policy
type Policy struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	Rules       []string `json:"rules"`
}

// PolicyListData represents policy list data with counts
type PolicyListData struct {
	TotalPolicies    int      `json:"total_policies"`
	ActivePolicies   int      `json:"active_policies"`
	InactivePolicies int      `json:"inactive_policies"`
	List             []Policy `json:"list"`
}

// ResponsePolicyList represents a list of policies response
type ResponsePolicyList struct {
	Status string         `json:"status"`
	Data   PolicyListData `json:"data"`
	Error  *string        `json:"error"`
}

// ResponsePolicySingle represents a single policy response
type ResponsePolicySingle struct {
	Status string  `json:"status"`
	Data   Policy  `json:"data"`
	Error  *string `json:"error"`
}

var mockPolicies = []Policy{
	{ID: "1", Name: "Block Ads", Description: "Block advertising domains", Active: true, Rules: []string{"*.ads.com", "*.doubleclick.net"}},
	{ID: "2", Name: "Block Malware", Description: "Block malicious domains", Active: true, Rules: []string{"*.malware.com"}},
	{ID: "3", Name: "Custom Policy", Description: "Custom blocking rules", Active: false, Rules: []string{}},
}

// ListPolicies handles GET /policies
func (h *APIHandler) ListPolicies(c *gin.Context) {
	// TODO: Implement logic to fetch policies from database
	activeCount := 0
	for _, p := range mockPolicies {
		if p.Active {
			activeCount++
		}
	}
	c.JSON(http.StatusOK, ResponsePolicyList{
		Status: "success",
		Data: PolicyListData{
			TotalPolicies:    len(mockPolicies),
			ActivePolicies:   activeCount,
			InactivePolicies: len(mockPolicies) - activeCount,
			List:             mockPolicies,
		},
		Error: nil,
	})
}

// GetPolicy handles GET /policies/:id
func (h *APIHandler) GetPolicy(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to fetch policy from database
	for _, policy := range mockPolicies {
		if policy.ID == id {
			c.JSON(http.StatusOK, ResponsePolicySingle{
				Status: "success",
				Data:   policy,
				Error:  nil,
			})
			return
		}
	}
	errMsg := "policy not found"
	c.JSON(http.StatusNotFound, ResponsePolicySingle{
		Status: "error",
		Data:   Policy{},
		Error:  &errMsg,
	})
}

// CreatePolicy handles POST /policies
func (h *APIHandler) CreatePolicy(c *gin.Context) {
	var policy Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, ResponsePolicySingle{
			Status: "error",
			Data:   Policy{},
			Error:  &errMsg,
		})
		return
	}
	// TODO: Implement logic to create policy in database
	mockPolicies = append(mockPolicies, policy)
	c.JSON(http.StatusCreated, ResponsePolicySingle{
		Status: "success",
		Data:   policy,
		Error:  nil,
	})
}

// DeletePolicy handles DELETE /policies/:id
func (h *APIHandler) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to delete policy from database
	for i, policy := range mockPolicies {
		if policy.ID == id {
			mockPolicies = append(mockPolicies[:i], mockPolicies[i+1:]...)
			c.JSON(http.StatusOK, ResponseGeneric{
				Status: "success",
				Data:   map[string]interface{}{},
				Error:  nil,
			})
			return
		}
	}
	errMsg := "policy not found"
	c.JSON(http.StatusNotFound, ResponseGeneric{
		Status: "error",
		Data:   nil,
		Error:  &errMsg,
	})
}
