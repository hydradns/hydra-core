package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DnsRecord struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}

type ResponseDnsSingle struct {
	Status string     `json:"status"`
	Data   DnsRecord  `json:"data"`
	Error  string     `json:"error,omitempty"`
}

type ResponseDnsList struct {
	Status string      `json:"status"`
	Data   []DnsRecord `json:"data"`
	Error  string      `json:"error,omitempty"`
}

var dnsRecords = []DnsRecord{
	{ID: "1", Domain: "example.com", Type: "A", Value: "192.168.1.1", TTL: 3600},
	{ID: "2", Domain: "test.com", Type: "AAAA", Value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", TTL: 3600},
}

// ListDNSRecords handles GET /dns
func (h *APIHandler) ListDNSRecords(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseDnsList{
		Status: "success",
		Data:   dnsRecords,
	})
}

// CreateDNSRecord handles POST /dns
func (h *APIHandler) CreateDNSRecord(c *gin.Context) {
	var newRecord DnsRecord
	if err := c.ShouldBindJSON(&newRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dnsRecords = append(dnsRecords, newRecord)
	c.JSON(http.StatusCreated, ResponseDnsSingle{
		Status: "success",
		Data:   newRecord,
	})
}

// DeleteDNSRecord handles DELETE /dns/:id
func (h *APIHandler) DeleteDNSRecord(c *gin.Context) {
	id := c.Param("id")
	for i, record := range dnsRecords {
		if record.ID == id {
			dnsRecords = append(dnsRecords[:i], dnsRecords[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{}})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
}
