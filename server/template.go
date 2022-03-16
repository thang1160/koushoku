package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"koushoku/cache"
	. "koushoku/config"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RenderOptions struct {
	Cache  bool
	Data   map[string]any
	Name   string
	Status int
}

const htmlContentType = "text/html; charset=utf-8"

var (
	templates           *template.Template
	mu                  sync.Mutex
	ErrTemplateNotFound = errors.New("Template not found")
)

func LoadTemplates() {
	mu.Lock()
	defer mu.Unlock()

	var files []string
	err := filepath.Walk(filepath.Join(Config.Directories.Templates),
		func(path string, stat fs.FileInfo, err error) error {
			if err != nil || stat.IsDir() || !strings.HasSuffix(path, ".html") {
				return err
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Fatalln(err)
	}

	templates, err = template.New("").Funcs(helper).ParseFiles(files...)
	if err != nil {
		log.Fatalln(err)
	}
}

func parseTemplate(name string, data any) ([]byte, error) {
	if gin.Mode() == gin.DebugMode {
		LoadTemplates()
	}

	t := templates.Lookup(name)
	if t == nil {
		return nil, ErrTemplateNotFound
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func getTemplate(name, key string) ([]byte, bool) {
	var v any
	var err error

	if len(key) > 0 {
		v, err = cache.Templates.Get(fmt.Sprintf("%s:%s", name, key))
	} else {
		v, err = cache.Templates.Get(name)
	}

	if err != nil {
		return nil, false
	}
	return v.([]byte), true
}

func setTemplate(name, key string, data any) ([]byte, error) {
	buf, err := parseTemplate(name, data)
	if err != nil {
		return nil, err
	}

	if len(key) > 0 {
		cache.Templates.Set(fmt.Sprintf("%s:%s", name, key), buf, 0)
	} else {
		cache.Templates.Set(name, buf, 0)
	}
	return buf, nil
}

func renderTemplate(c *Context, opts *RenderOptions) {
	var buf []byte
	if opts.Cache {
		var ok bool
		if buf, ok = getTemplate(opts.Name, c.cacheKey()); !ok {
			var err error
			buf, err = setTemplate(opts.Name, c.cacheKey(), opts.Data)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
	} else {
		buf, _ = parseTemplate(opts.Name, opts.Data)
	}
	c.Data(opts.Status, htmlContentType, buf)
}
