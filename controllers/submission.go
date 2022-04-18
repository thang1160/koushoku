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
	submitTemplateName      = "submit.html"
	submissionsTemplateName = "submissions.html"
)

func Submit(c *server.Context) {
	if c.TryCache(submitTemplateName) {
		return
	}
	c.Cache(http.StatusOK, submitTemplateName)
}

type SubmitPayload struct {
	Name      string `form:"name"`
	Submitter string `form:"submitter"`
	Content   string `form:"content"`
}

func SubmitPost(c *server.Context) {
	payload := &SubmitPayload{}
	c.Bind(payload)

	_, err := services.CreateSubmission(payload.Name, payload.Submitter, payload.Content)
	if err != nil {
		c.SetData("error", err)
		c.HTML(http.StatusBadRequest, submitTemplateName)
		return
	}

	c.SetData("message", "Your submission has been submitted.")
	c.HTML(http.StatusOK, submitTemplateName)
}

func Submissions(c *server.Context) {
	if c.TryCache(submissionsTemplateName) {
		return
	}

	q := &SearchQueries{}
	c.BindQuery(q)

	page, _ := strconv.Atoi(c.Query("page"))
	opts := services.GetSubmissionsOptions{
		Limit:  listingLimit,
		Offset: listingLimit * (page - 1),
	}

	result := services.GetSubmissions(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Submissions: Page %d", page))
	} else {
		c.SetData("name", "Submissions")
	}

	c.SetData("data", result.Submissions)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(listingLimit)))
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, submissionsTemplateName)
}
