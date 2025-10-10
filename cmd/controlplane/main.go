// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"controlplane/routes"
	"controlplane/middlewares"
	"controlplane/config"
	"github.com/gin-gonic/gin"
)

func main(){
	r:= gin.Default()

	//Using middleware for example... Logger
	r.Use(middlewares.Logger())

	routes.RegisterRoutes(r)

	r.Run(config.GetPort())
}