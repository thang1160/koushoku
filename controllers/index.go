package controllers

import (
	"fmt"
	"koushoku/server"
	"koushoku/services"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	indexLimit         = 25
	indexTemplateName  = "index.html"
	searchTemplateName = "search.html"
)

func Index(c *server.Context) {
	if c.TryCache(indexTemplateName) {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetArchivesOptions{
		Limit:  indexLimit,
		Offset: indexLimit * (page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
		},
	}
	result := services.GetArchives(c.IsUohhhhhhhhh(), opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Browse: Page %d", page))
	} else {
		c.SetData("name", "Home")
	}

	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, indexTemplateName)
}

type SearchQueries struct {
	Search string `form:"q"`
	Page   int    `form:"page"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

func Search(c *server.Context) {
	if c.TryCache(searchTemplateName) {
		return
	}

	q := SearchQueries{}
	c.BindQuery(&q)

	opts := services.GetArchivesOptions{
		Limit:  indexLimit,
		Offset: indexLimit * (q.Page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circle,
			services.ArchiveRels.Magazine,
		},
	}

	if len(q.Search) > 0 {
		ok := c.IsUohhhhhhhhh()

		if services.IsArtistValid(q.Search) {
			opts.Artists = []string{q.Search}
		}

		if services.IsCircleValid(q.Search) {
			opts.Circle = q.Search
		}

		if services.IsParodyValid(q.Search) {
			opts.Parody = q.Search
		}

		if services.IsTagValid(q.Search, ok) {
			opts.Tags = []string{q.Search}
		} else {
			arr := strings.Split(q.Search, " ")
			if len(arr) > 1 {
				for _, v := range arr {
					if services.IsTagValid(v, ok) {
						opts.Tags = append(opts.Tags, v)
					}
				}
			}
		}

		if len(opts.Artists) == 0 && len(opts.Circle) == 0 &&
			len(opts.Magazine) == 0 && len(opts.Parody) == 0 &&
			len(opts.Tags) == 0 {
			opts.Path = q.Search
		}
	}

	result := services.GetArchives(c.IsUohhhhhhhhh(), opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("queries", q)
	hasQueries := len(q.Search) > 0
	c.SetData("hasQueries", hasQueries)

	if hasQueries {
		c.SetData("name", fmt.Sprintf("Search: %s", q.Search))
	} else {
		c.SetData("name", "Browse")
	}

	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(q.Page, totalPages))

	c.Cache(http.StatusOK, searchTemplateName)
}
