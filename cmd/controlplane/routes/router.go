package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lopster568/phantomDNS/cmd/controlplane/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck)
	r.GET("/", handlers.Root)
}
