// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"log"

	"github.com/lopster568/phantomDNS/cmd/controlplane/config"
	"github.com/lopster568/phantomDNS/cmd/controlplane/handlers"
	"github.com/lopster568/phantomDNS/cmd/controlplane/middlewares"
	"github.com/lopster568/phantomDNS/cmd/controlplane/routes"
	"github.com/lopster568/phantomDNS/internal/grpc/client"
	"github.com/lopster568/phantomDNS/internal/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize grpc client
	c, err := client.New("dataplane:50051")
	if err != nil {
		log.Fatalf("failed to connect to dataplane: %v", err)
	}
	defer c.Close()

	status, err := c.CheckHealth()
	if err != nil {
		log.Fatalf("health check failed: %v", err)
	}

	logger.Log.Infof("Dataplane Health: %s\n", status)

	apiHandler := handlers.NewAPIHandler()

	// Initialize Gin router
	r := gin.Default()
	r.Use(middlewares.Logger())

	// CORS middleware (development-friendly). See cmd/controlplane/middlewares/cors.go
	r.Use(middlewares.CORS())

	routes.RegisterRoutes(r, apiHandler)
	r.Run(config.GetPort())
}
