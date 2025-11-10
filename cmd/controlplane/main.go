// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"github.com/lopster568/phantomDNS/cmd/controlplane/config"
	"github.com/lopster568/phantomDNS/cmd/controlplane/handlers"
	"github.com/lopster568/phantomDNS/cmd/controlplane/middlewares"
	"github.com/lopster568/phantomDNS/cmd/controlplane/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	r.Use(middlewares.Logger())

	// Create the API handler with dependencies
	apiHandler := handlers.NewAPIHandler()

	routes.RegisterRoutes(r, apiHandler)
	r.Run(config.GetPort())
}
