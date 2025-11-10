package routes

import (
	"github.com/lopster568/phantomDNS/cmd/controlplane/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, apiHandler *handlers.APIHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", apiHandler.HealthCheck)
		api.GET("/", apiHandler.Root)

		dns := api.Group("/dns")
		{
			dns.GET("", apiHandler.ListDNSRecords)
			dns.POST("", apiHandler.CreateDNSRecord)
			dns.DELETE("/:id", apiHandler.DeleteDNSRecord)
		}

		policies := api.Group("/policies")
		{
			policies.GET("", apiHandler.ListPolicies)
			policies.POST("", apiHandler.CreatePolicy)
			policies.DELETE("/:id", apiHandler.DeletePolicy)
		}

		blocklists := api.Group("/blocklists")
		{
			blocklists.GET("", apiHandler.ListBlocklists)
			blocklists.POST("", apiHandler.AddToBlocklist)
		}

		api.GET("/analytics", apiHandler.GetAnalyticsSummary)

		system := api.Group("/system")
		{
			system.GET("/status", apiHandler.GetSystemStatus)
		}
	}
}
