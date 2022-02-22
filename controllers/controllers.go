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

	server.GET("/archive/:id/:slug", Archive)
	server.GET("/archive/:id/:slug/:pageNum", ReadArchive)

	server.GET("/data/:id/:pageNum", ServeArchiveFile)
	server.GET("/data/:id/:pageNum/*width", ServeArchiveFile)

	server.GET("/artists", Artists)
	server.GET("/artists/:slug", Artist)
	server.GET("/circles", Circles)
	server.GET("/circles/:slug", Circle)
	server.GET("/tags", Tags)
	server.GET("/tags/:slug", Tag)
	server.GET("/magazines", Magazines)
	server.GET("/magazines/:slug", Magazine)
}

func Stats(c *server.Context) {
	c.SetData("stats", services.GetStats())
	c.HTML(http.StatusOK, "stats.html")
}
