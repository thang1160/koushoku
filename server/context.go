package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	. "koushoku/config"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context

	sync.RWMutex
	MapData map[string]interface{}
}

func (c *Context) GetURL() string {
	u, _ := url.Parse(Config.Meta.BaseURL)
	u.Path = c.Request.URL.Path
	u.RawQuery = c.Request.URL.RawQuery
	return u.String()
}

func (c *Context) preHTML(code *int) {
	if err, ok := c.GetData("error"); ok {
		err := strings.ToLower(err.(error).Error())
		if strings.Contains(err, "does not exist") || strings.Contains(err, "not found") {
			*code = http.StatusNotFound
		}
	}

	c.SetData("status", *code)
	c.SetData("statusText", http.StatusText(*code))

	if v, ok := c.MapData["name"]; !ok || len(v.(string)) == 0 {
		c.SetData("name", http.StatusText(*code))
	}

	c.SetData("title", Config.Meta.Title)
	c.SetData("description", Config.Meta.Description)
	c.SetData("baseURL", Config.Meta.BaseURL)
	c.SetData("language", Config.Meta.Language)

	c.SetData("url", c.GetURL())
	c.SetData("query", c.Request.URL.Query())
}

func (c *Context) HTML(code int, name string) {
	c.preHTML(&code)
	renderTemplate(c, false, &RenderOptions{
		Status: code,
		Name:   name,
		Data:   c.MapData,
	})
}

func (c *Context) Cache(code int, name string) {
	if gin.Mode() == gin.DebugMode {
		c.HTML(code, name)
	} else {
		c.preHTML(&code)
		renderTemplate(c, true, &RenderOptions{
			Status: code,
			Name:   name,
			Data:   c.MapData,
		})
	}
}

var getSecretKey func() string
var getSecretExpectedValue func() string

func (c *Context) IsUohhhhhhhhh() bool {
	if getSecretKey == nil || getSecretExpectedValue == nil {
		return false
	}
	v, _ := c.Cookie(getSecretKey())
	return v == getSecretExpectedValue()
}

func (c *Context) cacheKey() string {
	return fmt.Sprintf("%s%v", c.GetURL(), c.IsUohhhhhhhhh())
}

func (c *Context) IsCached(name string) bool {
	_, ok := getTemplate(name, c.cacheKey())
	return ok
}

func (c *Context) TryCache(name string) bool {
	if c.IsCached(name) {
		c.Cache(http.StatusOK, name)
		return true
	}
	return false
}

func (c *Context) ErrorJSON(code int, message string, err error) {
	c.JSON(code, gin.H{
		"error": gin.H{
			"message": message,
			"cause":   err.Error(),
		},
	})
}

func (c *Context) GetData(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	v, exists := c.MapData[key]
	return v, exists
}

func (c *Context) SetData(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	if c.MapData == nil {
		c.MapData = make(map[string]interface{})
	}
	c.MapData[key] = value
}

func (c *Context) SetCookie(name, value string, expires *time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   c.Request.TLS != nil || strings.HasPrefix(Config.Meta.BaseURL, "https"),
		HttpOnly: true,
	}

	if expires == nil {
		cookie.MaxAge = -1
	} else {
		cookie.Expires = *expires
	}

	http.SetCookie(c.Writer, cookie)
}

func (c *Context) ParamInt(name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}

func (c *Context) ParamInt64(name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}
