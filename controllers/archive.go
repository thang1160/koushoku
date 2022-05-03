package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	. "koushoku/config"

	"koushoku/server"
	"koushoku/services"
)

const (
	archiveTemplateName = "archive.html"
	readerTemplateName  = "reader.html"
)

func Archive(c *server.Context) {
	if c.TryCache(archiveTemplateName) {
		return
	}

	id, err := c.ParamInt64("id")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	opts := services.GetArchiveOptions{
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circles,
			services.ArchiveRels.Magazines,
			services.ArchiveRels.Parodies,
			services.ArchiveRels.Tags,
			services.ArchiveRels.Submission,
		},
	}
	result := services.GetArchive(id, opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	if (result.Archive.RedirectId > 0) && (result.Archive.RedirectId != id) {
		c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d", result.Archive.RedirectId))
		return
	}

	slug := c.Param("slug")
	if !strings.EqualFold(slug, result.Archive.Slug) {
		c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s", result.Archive.ID, result.Archive.Slug))
		return
	}

	c.SetData("archive", result.Archive)
	c.Cache(http.StatusOK, archiveTemplateName)
}

func Download(c *server.Context) {
	id, err := c.ParamInt64("id")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	fp, err := services.GetArchiveSymlink(int(id))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	} else if len(fp) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.FileAttachment(fp, filepath.Base(fp))
}

func Read(c *server.Context) {
	if c.TryCache(readerTemplateName) {
		return
	}

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
	c.Cache(http.StatusOK, readerTemplateName)
}

func createThumbnail(c *server.Context, f io.Reader, fp string, width int) (ok bool) {
	tmp, err := os.CreateTemp("", "tmp-")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if _, err := io.Copy(tmp, f); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	opts := services.ResizeOptions{Width: width, Height: width * 3 / 2}
	if err := services.ResizeImage(tmp.Name(), fp, opts); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	return true
}

func ServePage(c *server.Context) {
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
	index := pageNum - 1

	str := strings.TrimPrefix(c.Param("width"), "/")
	width, _ := strconv.Atoi(strings.TrimSuffix(str, filepath.Ext(str)))

	path, err := services.GetArchiveSymlink(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	} else if len(path) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	var fp string
	if (pageNum == 1 && (width == 288 || width == 896)) || width == 320 {
		fp = filepath.Join(Config.Directories.Thumbnails, fmt.Sprintf("%d-%d.%d.webp", id, pageNum, width))
		if _, err := os.Stat(fp); err == nil {
			c.ServeFile(fp)
			return
		}
	}

	zf, err := zip.OpenReader(path)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer zf.Close()

	var files []*zip.File
	for _, f := range zf.File {
		stat := f.FileInfo()
		name := stat.Name()

		if stat.IsDir() || !services.IsImage(name) {
			continue
		}

		files = append(files, f)
	}

	if index > len(files) {
		c.Status(http.StatusNotFound)
		return
	}

	sort.SliceStable(files, func(i, j int) bool {
		return services.GetPageNum(filepath.Base(files[i].Name)) < services.GetPageNum(filepath.Base(files[j].Name))
	})

	file := files[index]
	stat := file.FileInfo()

	f, err := file.Open()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if len(fp) > 0 {
		if createThumbnail(c, f, fp, width) {
			c.ServeFile(fp)
		}
	} else {
		buf, err := io.ReadAll(f)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.ServeData(stat, bytes.NewReader(buf))
	}
}
