package controllers

import "koushoku/server"

func Init() {
	server.GET("/", Index)
	server.GET("/search", Search)

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
