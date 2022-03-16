package main

import (
	"net/http"

	. "koushoku/config"

	"koushoku/controllers"
	"koushoku/server"
)

func main() {
	server.Init()

	server.GET("/archive/:id/:slug/download", controllers.Download)
	server.GET("/data/:id/:pageNum", controllers.ServePage)
	server.GET("/data/:id/:pageNum/*width", controllers.ServePage)

	server.NoRoute(func(c *server.Context) {
		c.Redirect(http.StatusFound, Config.Meta.BaseURL)
	})
	server.Start(Config.Server.DataPort)
}
