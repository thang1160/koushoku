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

func Archive(c *server.Context) {
	id, err := c.ParamInt64("id")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	if strings.EqualFold(c.Query("download"), "true") {
		services.ServeArchive(id, c.Writer, c.Request)
		return
	}

	opts := services.GetArchiveOptions{
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
			services.ArchiveRels.Parody,
			services.ArchiveRels.Tags,
		},
		IsUohhhhhhhhh: c.IsUohhhhhhhhh(),
	}
	result := services.GetArchive(id, opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("archive", result.Archive)
	c.HTML(http.StatusOK, "archive.html")
}
