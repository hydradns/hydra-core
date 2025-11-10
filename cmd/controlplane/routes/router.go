package routes

import(
    "github.com/gin-gonic/gin"
	"controlplane/handlers"
)

func RegisterRoutes(r *gin.Engine, apiHandler *handlers.APIHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", apiHandler.HealthCheck)
		api.GET("/", apiHandler.Root)
	}
}