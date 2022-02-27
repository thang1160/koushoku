package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"koushoku/server"
	"koushoku/services"
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
			services.ArchiveRels.Circles,
			services.ArchiveRels.Magazines,
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

	q.Search = strings.TrimSpace(q.Search)

	opts := services.GetArchivesOptions{
		Limit:  indexLimit,
		Offset: indexLimit * (q.Page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circles,
			services.ArchiveRels.Magazines,
		},
	}

	if len(q.Search) > 0 {
		if services.IsArtistValid(q.Search) {
			opts.Artists = []string{q.Search}
		}

		if services.IsCircleValid(q.Search) {
			opts.Circles = []string{q.Search}
		}

		if services.IsParodyValid(q.Search) {
			opts.Parodies = []string{q.Search}
		}

		if services.IsTagValid(q.Search) {
			opts.Tags = []string{q.Search}
		} else {
			arr := strings.Split(q.Search, " ")
			if len(arr) > 1 {
				for _, v := range arr {
					if services.IsTagValid(v) {
						opts.Tags = append(opts.Tags, v)
					}
				}
			}
		}

		if len(opts.Artists) == 0 && len(opts.Circles) == 0 &&
			len(opts.Magazines) == 0 && len(opts.Parodies) == 0 &&
			len(opts.Tags) == 0 {
			opts.Path = q.Search
		}
	}

	result := services.GetArchives(opts)
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

	if len(result.Archives) > 0 {
		c.Cache(http.StatusOK, searchTemplateName)
	} else {
		c.Cache(http.StatusNotFound, searchTemplateName)
	}
}
