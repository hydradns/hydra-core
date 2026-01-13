package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lopster568/phantomDNS/cmd/controlplane/handlers"
)

func RegisterRoutes(r *gin.Engine, apiHandler *handlers.APIHandler) {
	api := r.Group("/api/v1")
	r.GET("/health", apiHandler.HealthCheck)
	r.GET("/", apiHandler.Root)
	{
		// Dashboard endpoints
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/summary", apiHandler.GetDashboardSummary)
		}

		// DNS Engine endpoints
		dns := api.Group("/dns")
		{
			dns.GET("/engine", apiHandler.GetDnsEngineStatus)
			dns.POST("/engine", apiHandler.ToggleDnsEngine)
			dns.GET("/resolvers", apiHandler.ListResolvers)
			dns.GET("/metrics", apiHandler.GetDnsMetrics)
			dns.POST("/resolvers", apiHandler.AddResolver)
			dns.DELETE("/resolvers/:id", apiHandler.DeleteResolver)
		}

		// Policies endpoints
		policies := api.Group("/policies")
		{
			policies.GET("", apiHandler.ListPolicies)
			policies.POST("", apiHandler.CreatePolicy)
			policies.GET("/:id", apiHandler.GetPolicy)
			policies.DELETE("/:id", apiHandler.DeletePolicy)
		}

		// Blocklists endpoints
		blocklists := api.Group("/blocklists")
		{
			blocklists.GET("", apiHandler.ListBlocklists)
			blocklists.POST("", apiHandler.CreateBlocklist)
			blocklists.GET("/:id", apiHandler.GetBlocklist)
			blocklists.DELETE("/:id", apiHandler.DeleteBlocklist)
		}

		// Analytics endpoints
		analytics := api.Group("/analytics")
		{
			analytics.GET("/summary", apiHandler.GetAnalyticsSummary)
			analytics.GET("/audits", apiHandler.GetAuditLogs)
		}
	}
}
