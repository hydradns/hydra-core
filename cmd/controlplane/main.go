// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"log"

	"github.com/lopster568/phantomDNS/cmd/controlplane/handlers"
	"github.com/lopster568/phantomDNS/cmd/controlplane/middlewares"
	"github.com/lopster568/phantomDNS/cmd/controlplane/routes"
	"github.com/lopster568/phantomDNS/internal/config"
	client "github.com/lopster568/phantomDNS/internal/grpc/controlplane"
	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/storage/db"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db.InitDB("/app/data/phantomdns.db")
	repos := repositories.NewStore(db.DB)

	// Initialize grpc client
	c, err := client.New(config.DefaultConfig.DataPlane.GRPCServer.ListenAddr)
	if err != nil {
		log.Fatalf("failed to connect to dataplane: %v", err)
	}
	defer c.Close()

	status, err := c.GetStatus()
	if err != nil {
		log.Fatalf("health check failed: %v", err)
	}

	logger.Log.Infof("DNS Engine Status: %s\n", status)

	apiHandler := handlers.NewAPIHandler(*repos, c)

	// Initialize Gin router
	r := gin.Default()
	r.Use(middlewares.Logger())

	// CORS middleware (development-friendly). See cmd/controlplane/middlewares/cors.go
	r.Use(middlewares.CORS())

	routes.RegisterRoutes(r, apiHandler)
	r.Run(config.DefaultConfig.ControlPlane.ListenAddr)
}
