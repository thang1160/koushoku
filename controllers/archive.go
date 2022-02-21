package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"koushoku/server"
	"koushoku/services"
)

func ServeArchiveFile(c *server.Context) {
	id, err := c.ParamInt("id")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	pageNum := services.GetPageNum(c.Param("pageNum"))
	if pageNum <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	str := strings.TrimPrefix(c.Param("width"), "/")
	width, _ := strconv.Atoi(strings.TrimSuffix(str, filepath.Ext(str)))

	services.ServeArchiveFile(id, pageNum-1, width, c.Writer, c.Request)
}
