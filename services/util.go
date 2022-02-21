package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	. "koushoku/config"

	"koushoku/modext"
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

func FormatArchive(archive *modext.Archive) string {
	var s string
	if archive.Circle != nil {
		if len(archive.Artists) == 1 {
			s = fmt.Sprintf("[%s (%s)] ", archive.Circle.Name, archive.Artists[0].Name)
		} else {
			s = fmt.Sprintf("[%s] ", archive.Circle.Name)
		}
	} else if len(archive.Artists) > 0 && len(archive.Artists) < 3 {
		s += "["
		for i, artist := range archive.Artists {
			if i > 0 {
				s += ", "
			}
			s += artist.Name
		}
		s += "] "
	}
	s += archive.Title
	if archive.Magazine != nil {
		s += fmt.Sprintf(" [%s]", archive.Magazine.Name)
	}
	return s
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

func stringsContains(slice []string, search string) bool {
	for _, str := range slice {
		if str == search {
			return true
		}
	}
	return false
}

func CreateArchiveSymlink(archive *modext.Archive) error {
	if archive == nil {
		return nil
	}

	symlink := filepath.Join(Config.Directories.Symlinks, strconv.Itoa(int(archive.ID)))
	return os.Symlink(archive.Path, symlink)
}

func GetArchiveSymlink(id int) (string, error) {
	symlink := filepath.Join(Config.Directories.Symlinks, strconv.Itoa(id))
	return os.Readlink(symlink)
}

func PurgeArchiveSymlinks() {
	if err := os.RemoveAll(Config.Directories.Symlinks); err != nil {
		log.Fatalln(err)
	}
}

type ResizeOptions struct {
	Width  int
	Height int
	Crop   bool
}

var resizer struct {
	Map map[string]*sync.Mutex
	sync.Mutex
	sync.Once
}

func init() {
	resizer.Map = make(map[string]*sync.Mutex)
}

func resizeImage(filepath, outputPath string, o ResizeOptions) error {
	resizer.Lock()
	mu, ok := resizer.Map[outputPath]
	if !ok {
		mu = &sync.Mutex{}
		resizer.Map[outputPath] = mu
	}
	resizer.Unlock()

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
