// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"controlplane/routes"
	"controlplane/middlewares"
	"controlplane/config"
	"github.com/gin-gonic/gin"
)

func main(){
		// Initialize Gin router
	r := gin.Default()
	r.Use(middlewares.Logger())

	// Create the API handler with dependencies
	apiHandler := handlers.NewAPIHandler(c)

	routes.RegisterRoutes(r, apiHandler)
	r.Run(config.GetPort())
}