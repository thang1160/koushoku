package controllers

import (
	"net/http"

	"koushoku/server"
	"koushoku/services"
)

const (
	statsTemplateName = "stats.html"
)

func Stats(c *server.Context) {
	if c.TryCache(statsTemplateName) {
		return
	}

	c.SetData("stats", services.GetStats())
	c.Cache(http.StatusOK, statsTemplateName)
}
