package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"koushoku/server"
	"koushoku/services"
)

const (
	listingLimit    = 200
	listingTemplate = "list.html"
)

func Artists(c *server.Context) {
	if c.TryCache(listingTemplate) {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetArtistsOptions{
		Limit:  listingLimit,
		Offset: listingLimit * (page - 1),
	}

	result := services.GetArtists(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Artists: Page %d", page))
	} else {
		c.SetData("name", "Artists")
	}

	c.SetData("taxonomy", "artists")
	c.SetData("taxonomyTitle", "Artists")

	c.SetData("data", result.Artists)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(listingLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, listingTemplate)
}

func Circles(c *server.Context) {
	if c.TryCache(listingTemplate) {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetCirclesOptions{
		Limit:  listingLimit,
		Offset: listingLimit * (page - 1),
	}

	result := services.GetCircles(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Circles: Page %d", page))
	} else {
		c.SetData("name", "Circles")
	}

	c.SetData("taxonomy", "circles")
	c.SetData("taxonomyTitle", "Circles")

	c.SetData("data", result.Circles)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(listingLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, listingTemplate)
}

func Magazines(c *server.Context) {
	if c.TryCache(listingTemplate) {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetMagazinesOptions{
		Limit:  listingLimit,
		Offset: listingLimit * (page - 1),
	}

	result := services.GetMagazines(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Magazines: Page %d", page))
	} else {
		c.SetData("name", "Magazines")
	}

	c.SetData("taxonomy", "magazines")
	c.SetData("taxonomyTitle", "Magazines")

	c.SetData("data", result.Magazines)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(listingLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, listingTemplate)
}

func Parodies(c *server.Context) {
	if c.TryCache(listingTemplate) {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetParodiesOptions{
		Limit:  listingLimit,
		Offset: listingLimit * (page - 1),
	}

	result := services.GetParodies(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Parodies: Page %d", page))
	} else {
		c.SetData("name", "Parodies")
	}

	c.SetData("taxonomy", "parodies")
	c.SetData("taxonomyTitle", "Parodies")

	c.SetData("data", result.Parodies)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(listingLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, listingTemplate)
}

func Tags(c *server.Context) {
	if c.TryCache(listingTemplate) {
		return
	}

	result := services.GetTags(services.GetTagsOptions{})
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("name", "Tags")
	c.SetData("taxonomy", "tags")
	c.SetData("taxonomyTitle", "Tags")

	c.SetData("data", result.Tags)
	c.SetData("total", result.Total)

	c.Cache(http.StatusOK, listingTemplate)
}
