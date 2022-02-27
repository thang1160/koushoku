package controllers

import (
	"net/http"

	"koushoku/server"
	"koushoku/services"
)

func Init() {
	server.GET("/", Index)
	server.GET("/search", Search)
	server.GET("/stats",
		server.WithName("Stats"),
		Stats)

	server.GET("/archive/:id", Archive)
	server.GET("/archive/:id/:slug", Archive)
	server.GET("/archive/:id/:slug/:pageNum", Read)

	server.GET("/data/:id/:pageNum", ServePage)
	server.GET("/data/:id/:pageNum/*width", ServePage)

	server.GET("/artists", Artists)
	server.GET("/artists/:slug", Artist)
	server.GET("/circles", Circles)
	server.GET("/circles/:slug", Circle)
	server.GET("/magazines", Magazines)
	server.GET("/magazines/:slug", Magazine)
	server.GET("/parodies", Parodies)
	server.GET("/parodies/:slug", Parody)
	server.GET("/tags", Tags)
	server.GET("/tags/:slug", Tag)
}

const statsTemplateName = "stats.html"

func Stats(c *server.Context) {
	if c.TryCache(statsTemplateName) {
		return
	}

	c.SetData("stats", services.GetStats())
	c.Cache(http.StatusOK, statsTemplateName)
}
