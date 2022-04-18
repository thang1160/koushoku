package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	. "koushoku/config"

	"koushoku/cache"
	"koushoku/controllers"
	"koushoku/database"
	"koushoku/server"
	"koushoku/services"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	cache.Init()

	if err := services.AnalyzeStats(); err != nil {
		return
	}
	server.Init()

	assets := server.Group("/")
	assets.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=300")
	})

	assets.Static("/js", filepath.Join(Config.Directories.Root, "assets/js"))
	assets.Static("/css", filepath.Join(Config.Directories.Root, "assets/css"))
	assets.Static("/fonts", filepath.Join(Config.Directories.Root, "assets/fonts"))

	server.GET("/archive/:id/:slug/download", func(c *server.Context) {
		c.Redirect(http.StatusFound,
			fmt.Sprintf("%s/archive/%s/%s/download",
				Config.Meta.DataBaseURL,
				c.Param("id"),
				c.Param("slug")))
	})
	server.GET("/data/:id/:pageNum", func(c *server.Context) {
		c.Redirect(http.StatusFound,
			fmt.Sprintf("%s/data/%s/%s",
				Config.Meta.DataBaseURL,
				c.Param("id"),
				c.Param("pageNum")))
	})
	server.GET("/data/:id/:pageNum/*width", func(c *server.Context) {
		c.Redirect(http.StatusFound,
			fmt.Sprintf("%s/data/%s/%s/%s",
				Config.Meta.DataBaseURL,
				c.Param("id"),
				c.Param("pageNum"),
				c.Param("width")))
	})

	// assets.StaticFile("/serviceWorker.js", filepath.Join(Config.Directories.Root, "assets/js/serviceWorker.js"))
	assets.StaticFile("/app.webmanifest", filepath.Join(Config.Directories.Root, "app.webmanifest"))
	assets.StaticFile("/cover.jpg", filepath.Join(Config.Directories.Root, "cover.jpg"))

	assets.StaticFile("/robots.txt", filepath.Join(Config.Directories.Root, "robots.txt"))
	assets.StaticFile("/updates.txt", filepath.Join(Config.Directories.Root, "updates.txt"))

	assets.StaticFile("/favicon.ico", filepath.Join(Config.Directories.Root, "favicon.ico"))
	assets.StaticFile("/favicon-16x16.png", filepath.Join(Config.Directories.Root, "favicon-16x16.png"))
	assets.StaticFile("/favicon-32x32.png", filepath.Join(Config.Directories.Root, "favicon-32x32.png"))
	assets.StaticFile("/apple-touch-icon.png", filepath.Join(Config.Directories.Root, "apple-touch-icon.png"))
	assets.StaticFile("/android-chrome-192x192.png", filepath.Join(Config.Directories.Root, "android-chrome-192x192.png"))
	assets.StaticFile("/android-chrome-512x512.png", filepath.Join(Config.Directories.Root, "android-chrome-512x512.png"))

	server.GET("/", controllers.Index)
	server.GET("/search", controllers.Search)
	server.GET("/stats",
		server.WithName("Stats"),
		controllers.Stats)

	server.GET("/archive/:id", controllers.Archive)
	server.GET("/archive/:id/:slug", controllers.Archive)
	server.GET("/archive/:id/:slug/:pageNum", controllers.Read)

	server.GET("/artists", controllers.Artists)
	server.GET("/artists/:slug", controllers.Artist)

	server.GET("/circles", controllers.Circles)
	server.GET("/circles/:slug", controllers.Circle)

	server.GET("/magazines", controllers.Magazines)
	server.GET("/magazines/:slug", controllers.Magazine)

	server.GET("/parodies", controllers.Parodies)
	server.GET("/parodies/:slug", controllers.Parody)

	server.GET("/tags", controllers.Tags)
	server.GET("/tags/:slug", controllers.Tag)

	server.GET("/submit", server.WithName("Submit"), controllers.Submit)
	server.POST("/submit",
		server.WithName("Submit"),
		server.WithRateLimit("Submit?", "10-D"),
		controllers.SubmitPost)
	server.GET("/submissions", controllers.Submissions)

	server.GET("/sitemap.xml", controllers.Sitemap)

	server.NoRoute(func(c *server.Context) {
		c.HTML(http.StatusNotFound, "error.html")
	})
	server.Start(Config.Server.WebPort)
}
