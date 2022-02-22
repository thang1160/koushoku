package controllers

import (
	"fmt"
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

func ReadArchive(c *server.Context) {
	id, err := c.ParamInt64("id")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	pageNum, err := c.ParamInt("pageNum")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	opts := services.GetArchiveOptions{}
	result := services.GetArchive(id, opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	slug := c.Param("slug")
	if !strings.EqualFold(slug, result.Archive.Slug) {
		if pageNum <= 0 || int16(pageNum) > result.Archive.Pages {
			c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s/1", result.Archive.ID, result.Archive.Slug))
		} else {
			c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s/%d", result.Archive.ID, result.Archive.Slug, pageNum))
		}
		return
	}

	if pageNum <= 0 || int16(pageNum) > result.Archive.Pages {
		c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s/1", id, result.Archive.Slug))
		return
	}

	c.SetData("archive", result.Archive)
	c.SetData("pageNum", pageNum)

	c.HTML(http.StatusOK, "reader.html")
}
