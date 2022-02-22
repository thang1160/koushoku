package controllers

import (
	"fmt"
	"koushoku/server"
	"koushoku/services"
	"math"
	"net/http"
	"strconv"
)

const (
	taxonomyTemplate = "taxonomy.html"
)

func Artist(c *server.Context) {
	if c.TryCache(taxonomyTemplate) {
		return
	}

	artist, err := services.GetArtist(c.Param("slug"))
	if err != nil {
		c.SetData("error", err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetArchivesOptions{
		Artists: []string{artist.Name},
		Limit:   indexLimit,
		Offset:  indexLimit * (page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
		},
	}

	result := services.GetArchives(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("%s: Page %d", artist.Name, page))
	} else {
		c.SetData("name", artist.Name)
	}

	c.SetData("taxonomy", artist.Name)
	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, taxonomyTemplate)
}

func Circle(c *server.Context) {
	if c.TryCache(taxonomyTemplate) {
		return
	}

	circle, err := services.GetCircle(c.Param("slug"))
	if err != nil {
		c.SetData("error", err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetArchivesOptions{
		Circle: circle.Name,
		Limit:  indexLimit,
		Offset: indexLimit * (page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
		},
	}

	result := services.GetArchives(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("%s: Page %d", circle.Name, page))
	} else {
		c.SetData("name", circle.Name)
	}

	c.SetData("taxonomy", circle.Name)
	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, taxonomyTemplate)
}

func Tag(c *server.Context) {
	if c.TryCache(taxonomyTemplate) {
		return
	}

	tag, err := services.GetTag(c.Param("slug"))
	if err != nil {
		c.SetData("error", err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetArchivesOptions{
		Tags:   []string{tag.Name},
		Limit:  indexLimit,
		Offset: indexLimit * (page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
		},
	}

	result := services.GetArchives(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("%s: Page %d", tag.Name, page))
	} else {
		c.SetData("name", tag.Name)
	}

	c.SetData("taxonomy", tag.Name)
	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, taxonomyTemplate)
}

func Magazine(c *server.Context) {
	if c.TryCache(taxonomyTemplate) {
		return
	}

	magazine, err := services.GetMagazine(c.Param("slug"))
	if err != nil {
		c.SetData("error", err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetArchivesOptions{
		Magazine: magazine.Name,
		Limit:    indexLimit,
		Offset:   indexLimit * (page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
		},
	}

	result := services.GetArchives(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("%s: Page %d", magazine.Name, page))
	} else {
		c.SetData("name", magazine.Name)
	}

	c.SetData("taxonomy", magazine.Name)
	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, taxonomyTemplate)
}
