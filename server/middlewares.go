package server

import "net/http"

func WithName(name string) Handler {
	return func(c *Context) {
		c.SetData("name", name)
		c.Next()
	}
}

func WithRedirect(relativePath string) Handler {
	return func(c *Context) {
		c.Redirect(http.StatusFound, relativePath)
	}
}
