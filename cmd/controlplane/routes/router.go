package routes

import(
    "github.com/gin-gonic/gin"
	"controlplane/handlers"
)

func RegisterRoutes(r *gin.Engine) {
    r.GET("/health", handlers.HealthCheck)
	r.GET("/", handlers.Root)
}