// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"github.com/lopster568/phantomDNS/cmd/controlplane/config"
	"github.com/lopster568/phantomDNS/cmd/controlplane/middlewares"
	"github.com/lopster568/phantomDNS/cmd/controlplane/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//Using middleware for example... Logger
	r.Use(middlewares.Logger())

	routes.RegisterRoutes(r)

	r.Run(config.GetPort())
}
