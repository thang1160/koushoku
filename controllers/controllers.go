package controllers

import (
	"net/http"

	"koushoku/server"
	"koushoku/services"
)

const (
	aboutTemplateName = "about.html"
	statsTemplateName = "stats.html"
)

func About (c *server.Context) {
	if c.TryCache(aboutTemplateName) {
		return;
	}

	c.Cache(http.StatusOK, aboutTemplateName)
}

func Stats(c *server.Context) {
	if c.TryCache(statsTemplateName) {
		return
	}

	c.SetData("stats", services.GetStats())
	c.Cache(http.StatusOK, statsTemplateName)
}
