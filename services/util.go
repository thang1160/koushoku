package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gosimple/slug"
)

type QueryMapCache struct {
	Map map[string]bool
	sync.RWMutex
}

type Pagination struct {
	CurrentPage int
	Pages       []int
	TotalPages  int
}

const maxPages = 10

func CreatePagination(currentPage, totalPages int) *Pagination {
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	pagination := &Pagination{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}

	var first, last int
	if totalPages <= maxPages {
		first = 1
		last = totalPages
	} else {
		min := int(math.Floor(float64(maxPages) / 2))
		max := int(math.Ceil(float64(maxPages)/2)) - 1
		if currentPage <= min {
			first = 1
			last = maxPages
		} else if currentPage+max >= totalPages {
			first = totalPages - maxPages + 1
			last = totalPages
		} else {
			first = currentPage - min
			last = currentPage + max
		}
	}

	pagination.Pages = make([]int, last-first+1)
	for i := 0; i < last+1-first; i++ {
		pagination.Pages[i] = first + i
	}

	return pagination
}

func FormatBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func FileName(path string) string {
	return strings.TrimRight(filepath.Base(path), filepath.Ext(path))
}

var rgx = regexp.MustCompile("[0-9]+")

func GetPageNum(fileName string) int {
	fileName = strings.TrimLeft(fileName, "0")
	n, _ := strconv.Atoi(rgx.FindString(fileName))
	return n
}

func JoinURL(base string, paths ...string) string {
	u, _ := url.Parse(base)
	for _, path := range paths {
		u.Path = filepath.Join(u.Path, strings.TrimLeft(strings.TrimRight(path, "/"), "/"))
	}
	return u.String()
}

func makeCacheKey(v interface{}) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

var slugCache struct {
	Map map[string]string
	sync.RWMutex
	sync.Once
}

func slugify(s string) string {
	slugCache.Once.Do(func() {
		slugCache.Map = make(map[string]string)
	})

	s = strings.ToLower(s)

	slugCache.RLock()
	v, ok := slugCache.Map[s]
	slugCache.RUnlock()

	if !ok {
		v = slug.Make(s)

		slugCache.Lock()
		slugCache.Map[s] = v
		slugCache.Unlock()
	}

	return v
}

func stringsContains(slice []string, search string) bool {
	for _, str := range slice {
		if str == search {
			return true
		}
	}
	return false
}

func pluralize(str string) string {
	if strings.HasSuffix(str, "ss") {
		return str + "es"
	} else if strings.HasSuffix(str, "y") {
		return strings.TrimSuffix(str, "y") + "ies"
	}
	return str + "s"
}

type ResizeOptions struct {
	Width  int
	Height int
	Crop   bool
}

var resizer struct {
	Map   map[string]*sync.Mutex
	Queue chan bool
	sync.RWMutex
	sync.Once
}

func init() {
	resizer.Map = make(map[string]*sync.Mutex)
	resizer.Queue = make(chan bool, 10)
}

func ResizeImage(filepath, outputPath string, o ResizeOptions) error {
	resizer.RLock()
	mu, ok := resizer.Map[outputPath]
	resizer.RUnlock()

	if !ok {
		mu = &sync.Mutex{}
		resizer.Lock()
		resizer.Map[outputPath] = mu
		resizer.Unlock()
	}

	mu.Lock()
	defer func() {
		mu.Unlock()

		resizer.Lock()
		delete(resizer.Map, outputPath)
		resizer.Unlock()
	}()

	if ok {
		return nil
	}

	resizer.Queue <- true
	defer func() {
		<-resizer.Queue
	}()

	w := strconv.Itoa(o.Width)
	h := strconv.Itoa(o.Height)
	crop := strconv.FormatBool(o.Crop)

	buf, err := runCommand("./resizer", filepath, w, h, crop)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0755)
}

func runCommand(path string, args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()

	var buf bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if err := stderr.String(); len(err) > 0 {
		return nil, errors.New(err)
	}

	return &buf, nil
}
