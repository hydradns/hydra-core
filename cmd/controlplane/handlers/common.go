package handlers

import "github.com/gin-gonic/gin"

func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *APIHandler) Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to PhantomDNS Control Plane API",
	})
}
