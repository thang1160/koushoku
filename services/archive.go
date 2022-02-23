package services

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	. "koushoku/cache"
	. "koushoku/config"

	"koushoku/errs"
	"koushoku/models"
	"koushoku/modext"

	"github.com/gosimple/slug"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	ArchiveCols = models.ArchiveColumns
	ArchiveRels = models.ArchiveRels
)

func checkArchiveFileThumbnail(id, index, width int, w http.ResponseWriter, r *http.Request) (string, bool) {
	var fp string
	if width > 0 && width <= 1024 && width%128 == 0 {
		fp = filepath.Join(Config.Directories.Thumbnails, fmt.Sprintf("%d-%d.%d.jpg", id, index+1, width))
		if _, err := os.Stat(fp); err == nil {
			http.ServeFile(w, r, fp)
			return fp, true
		}
	}
	return fp, false
}

func createArchiveThumbnail(f io.Reader, fp string, width int, w http.ResponseWriter, r *http.Request) (ok bool) {
	tmp, err := os.CreateTemp("", "tmp-")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if _, err := io.Copy(tmp, f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	opts := ResizeOptions{
		Width:  width,
		Height: width * 3 / 2,
		Crop:   true,
	}

	if err := resizeImage(tmp.Name(), fp, opts); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return true
}

func ServeArchive(id int64, w http.ResponseWriter, r *http.Request) {
	path, err := GetArchiveSymlink(int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if len(path) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	http.ServeFile(w, r, path)
}

func ServeArchiveFile(id, index, width int, w http.ResponseWriter, r *http.Request) {
	fp, ok := checkArchiveFileThumbnail(id, index, width, w, r)
	if ok {
		return
	}

	path, err := GetArchiveSymlink(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if len(path) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	z, err := zip.OpenReader(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer z.Close()

	if index > len(z.File) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var file *zip.File
	var d fs.FileInfo

	for true {
		file = z.File[index]
		d = file.FileInfo()

		if d.IsDir() {
			index++
			continue
		}
		break
	}

	f, err := file.Open()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if len(fp) > 0 {
		if createArchiveThumbnail(f, fp, width, w, r) {
			http.ServeFile(w, r, fp)
		}
	} else {
		buf, err := io.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, d.Name(), d.ModTime(), bytes.NewReader(buf))
	}
}

func IndexArchives() {
	if _, err := os.Stat(Config.Directories.Symlinks); os.IsNotExist(err) {
		if err := os.MkdirAll(Config.Directories.Symlinks, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	var archives []*modext.Archive
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() ||
			strings.Contains(path, "/cover") || strings.Contains(path, "/doujin") ||
			strings.Contains(path, "/illustration") || strings.Contains(path, "/interview") ||
			strings.Contains(path, "/non-h") || strings.Contains(path, "/spread") ||
			strings.Contains(path, "/western") || !strings.HasSuffix(path, ".zip") {
			return err
		}

		fileName := strings.TrimRight(filepath.Base(path), ".zip")
		archive := &modext.Archive{Path: path}
		archive.FormatFromString(fileName)

		if len(archive.Title) > 0 {
			log.Println("Found archive", fileName)
			archives = append(archives, archive)
		}

		return nil
	}

	filepath.Walk(Config.Directories.Data, walkFn)
	for _, archive := range archives {
		log.Println("Indexing archive", filepath.Base(archive.Path))
		c, err := CreateArchive(archive)
		if c != nil && err == nil {
			CreateArchiveSymlink(c)
		}
	}
}

func getBlacklists(file string) (artists, titles map[string]bool, err error) {
	if d, err := os.Stat(file); os.IsNotExist(err) {
		return nil, nil, err
	} else if d.IsDir() {
		return nil, nil, errors.New("Input is a directory")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	artists = make(map[string]bool)
	titles = make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		arr := strings.Split(line, ":")
		if len(arr) != 2 {
			continue
		}

		t := strings.TrimSpace(arr[0])
		v := slug.Make(arr[1])

		if strings.EqualFold(t, "artist") {
			artists[v] = true
		} else if strings.EqualFold(t, "title") {
			titles[v] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return
}

func ModerateArchives(file string) {
	artists, titles, err := getBlacklists(file)
	if err != nil {
		log.Fatalln(err)
	}

	archives, err := models.Archives(Load(ArchiveRels.Artists)).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	for _, archive := range archives {
		_, remove := titles[archive.Slug]
		if archive.R != nil && len(archive.R.Artists) > 0 {
			for _, artist := range archive.R.Artists {
				if _, ok := artists[artist.Slug]; ok {
					artist.DeleteG()
					remove = true
				}
			}
		}
		if remove {
			log.Println("Removing archive", archive.Path)
			DeleteArchive(archive.ID)
		}
	}
}

func refreshArchiveRels(arc *models.Archive, archive *modext.Archive) error {
	if archive.Circle != nil {
		circle, err := CreateCircle(archive.Circle.Name)
		if err != nil {
			return err
		}
		if err := arc.SetCircleG(false, &models.Circle{
			ID:   circle.ID,
			Slug: circle.Slug,
			Name: circle.Name,
		}); err != nil {
			log.Println(err)
			return errs.ErrUnknown
		}
	}

	if archive.Magazine != nil {
		magazine, err := CreateMagazine(archive.Magazine.Name)
		if err != nil {
			return err
		}
		if err := arc.SetMagazineG(false, &models.Magazine{
			ID:   magazine.ID,
			Slug: magazine.Slug,
			Name: magazine.Name,
		}); err != nil {
			log.Println(err)
			return errs.ErrUnknown
		}
	}

	if archive.Parody != nil {
		parody, err := CreateParody(archive.Parody.Name)
		if err != nil {
			return err
		}
		if err := arc.SetParodyG(false, &models.Parody{
			ID:   parody.ID,
			Slug: parody.Slug,
			Name: parody.Name,
		}); err != nil {
			log.Println(err)
			return errs.ErrUnknown
		}
	}

	var artists []*models.Artist
	for _, artist := range archive.Artists {
		artist, err := CreateArtist(artist.Name)
		if err != nil {
			return err
		}
		artists = append(artists, &models.Artist{
			ID:   artist.ID,
			Slug: artist.Slug,
			Name: artist.Name,
		})
	}
	if err := arc.SetArtistsG(false, artists...); err != nil {
		log.Println(err)
		return errs.ErrUnknown
	}

	var tags []*models.Tag
	for _, tag := range archive.Tags {
		tag, err := CreateTag(tag.Name)
		if err != nil {
			return err
		}
		tags = append(tags, &models.Tag{
			ID:   tag.ID,
			Slug: tag.Slug,
			Name: tag.Name,
		})
	}
	if err := arc.SetTagsG(false, tags...); err != nil {
		log.Println(err)
		return errs.ErrUnknown
	}
	return nil
}

func CreateArchive(archive *modext.Archive) (*modext.Archive, error) {
	if archive == nil {
		return nil, nil
	} else if len(archive.Path) == 0 {
		return nil, errs.ErrArchivePathRequired
	}

	stat, err := os.Stat(archive.Path)
	if os.IsNotExist(err) {
		log.Println(err)
		return nil, errs.ErrArchiveNotFound
	}

	arc, err := models.Archives(Where("archive.path = ?", archive.Path)).OneG()
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, errs.ErrUnknown
	} else if arc != nil {
		return modext.NewArchive(arc), nil
	}

	arc = &models.Archive{
		Path: archive.Path,

		Title: archive.Title,
		Slug:  slug.Make(archive.Title),
		Pages: archive.Pages,
		Size:  FormatBytes(stat.Size()),
	}

	f, err := zip.OpenReader(arc.Path)
	if err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	for _, file := range f.File {
		if file.FileInfo().IsDir() {
			continue
		}
		arc.Pages++
	}

	if arc.Pages > 0 {
		d := f.File[0].FileInfo()
		arc.CreatedAt = d.ModTime()
		arc.UpdatedAt = arc.CreatedAt
	}
	f.Close()

	if err := arc.InsertG(boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	} else if err := refreshArchiveRels(arc, archive); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	// TODO: Purge cache
	return modext.NewArchive(arc), nil
}

func validateRels(rels []string) (result []string) {
	for _, v := range rels {
		if strings.EqualFold(v, ArchiveRels.Artists) {
			result = append(result, ArchiveRels.Artists)
		} else if strings.EqualFold(v, ArchiveRels.Circle) {
			result = append(result, ArchiveRels.Circle)
		} else if strings.EqualFold(v, ArchiveRels.Magazine) {
			result = append(result, ArchiveRels.Magazine)
		} else if strings.EqualFold(v, ArchiveRels.Parody) {
			result = append(result, ArchiveRels.Parody)
		} else if strings.EqualFold(v, ArchiveRels.Tags) {
			result = append(result, ArchiveRels.Tags)
		}
	}
	sort.Strings(result)
	return
}

type GetArchiveOptions struct {
	Preloads []string `form:"preload" json:"1,omitempty"`
}

type GetArchiveResult struct {
	Archive *modext.Archive `json:"archive,omitempty"`
	Err     error           `json:"error,omitempty"`
}

func GetArchive(id int64, opts GetArchiveOptions) (result *GetArchiveResult) {
	opts.Preloads = validateRels(opts.Preloads)

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(id, cacheKey); err == nil {
		return c.(*GetArchiveResult)
	}

	result = &GetArchiveResult{}
	defer func() {
		if result.Archive != nil || result.Err != nil {
			Cache.RemoveWithPrefix(id, cacheKey)
			Cache.SetWithPrefix(id, cacheKey, result, 0)
		}
	}()

	selectQueries := []QueryMod{
		Where("id = ?", id),
		And("published_at IS NOT NULL"),
	}

	for _, v := range opts.Preloads {
		selectQueries = append(selectQueries, Load(v))
	}

	archive, err := models.Archives(selectQueries...).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			result.Err = errs.ErrArchiveNotFound
			return
		}
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	result.Archive = modext.NewArchive(archive).LoadRels(archive)
	return
}

type GetArchivesOptions struct {
	Path     string `json:"0,omitempty"`
	Title    string `json:"1,omitempty"`
	Circle   string `json:"2,omitempty"`
	Magazine string `json:"3,omitempty"`
	Parody   string `json:"4,omitempty"`

	Artists []string `json:"5,omitempty"`
	Tags    []string `json:"6,omitempty"`

	Limit    int      `json:"7,omitempty"`
	Offset   int      `json:"8,omitempty"`
	Preloads []string `json:"9,omitempty"`
	Sort     string   `json:"10,omitempty"`
	Order    string   `json:"11,omitempty"`
}

func (o *GetArchivesOptions) validate() {
	o.Path = strings.ToLower(o.Path)
	o.Title = slug.Make(o.Title)
	o.Circle = slug.Make(o.Circle)
	o.Magazine = slug.Make(o.Magazine)
	o.Parody = slug.Make(o.Parody)

	for i, artist := range o.Artists {
		o.Artists[i] = slug.Make(artist)
	}
	sort.Strings(o.Artists)

	for i, tag := range o.Tags {
		o.Tags[i] = slug.Make(tag)
	}
	sort.Strings(o.Tags)

	if o.Limit <= 0 {
		o.Limit = 20
	} else if o.Limit > 100 {
		o.Limit = 100
	}

	if o.Offset < 0 {
		o.Offset = 0
	}

	o.Preloads = validateRels(o.Preloads)

	if strings.EqualFold(o.Sort, ArchiveCols.ID) {
		o.Sort = ArchiveCols.ID
	} else if strings.EqualFold(o.Sort, ArchiveCols.UpdatedAt) {
		o.Sort = ArchiveCols.UpdatedAt
	} else if strings.EqualFold(o.Sort, ArchiveCols.PublishedAt) {
		o.Sort = ArchiveCols.PublishedAt
	} else if strings.EqualFold(o.Sort, ArchiveCols.Title) {
		o.Sort = ArchiveCols.Title
	} else {
		o.Sort = ArchiveCols.CreatedAt
	}

	if strings.EqualFold(o.Order, "asc") {
		o.Order = "asc"
	} else {
		o.Order = "desc"
	}

}

func (o *GetArchivesOptions) toQueries(isOr bool) (selectQueries, countQueries []QueryMod) {
	countQueries = append(countQueries, Select("1"))

	var queries []string
	var args []interface{}

	if len(o.Path) > 0 {
		queries = append(queries, "archive.path ILIKE '%' || ? || '%'")
		args = append(args, o.Path)
	}

	if len(o.Title) > 0 {
		queries = append(queries, "archive.slug ILIKE '%' || ? || '%'")
		args = append(args, o.Title)
	}

	if len(o.Circle) > 0 {
		selectQueries = append(selectQueries,
			InnerJoin("circle ON circle.id = archive.circle_id"))
		if isOr {
			queries = append(queries, "circle.slug ILIKE '%' || ? || '%'")
		} else {
			queries = append(queries, "circle.slug = ?")
		}
		args = append(args, o.Circle)
	}

	if len(o.Magazine) > 0 {
		selectQueries = append(selectQueries,
			InnerJoin("magazine ON archive.magazine_id = magazine.id"))
		queries = append(queries, "magazine.slug = ?")
		args = append(args, o.Magazine)
	}

	if len(o.Parody) > 0 {
		selectQueries = append(selectQueries,
			InnerJoin("parody ON parody.id = archive.parody_id"))
		queries = append(queries, "parody.slug = ?")
		args = append(args, o.Parody)
	}

	if len(o.Artists) > 0 {
		selectQueries = append(selectQueries,
			InnerJoin("archive_artists ar ON ar.archive_id = archive.id"),
			InnerJoin("artist ON artist.id = ar.artist_id"))
		var q []string
		for _, tag := range o.Artists {
			if isOr {
				q = append(q, "artist.slug ILIKE '%' || ? || '%'")
			} else {
				q = append(q, "artist.slug = ?")
			}
			args = append(args, tag)
		}
		queries = append(queries, fmt.Sprintf("(%s)", strings.Join(q, " OR ")))
	}

	if len(o.Tags) > 0 {
		selectQueries = append(selectQueries,
			InnerJoin("archive_tags at ON at.archive_id = archive.id"),
			InnerJoin("tag ON tag.id = at.tag_id"))

		var q []string
		for _, tag := range o.Tags {
			q = append(q, "tag.slug = ?")
			args = append(args, tag)
		}
		queries = append(queries, fmt.Sprintf("(%s)", strings.Join(q, " OR ")))
	}

	if len(o.Artists) > 0 || len(o.Tags) > 0 {
		selectQueries = append(selectQueries, GroupBy("archive.id"))
	}

	if len(queries) > 0 {
		op := " AND "
		if isOr {
			op = " OR "
		}
		selectQueries = append(selectQueries,
			Where(strings.Join(queries, op), args...))
	}

	selectQueries = append(selectQueries, Where("archive.published_at IS NOT NULL"))
	countQueries = append(countQueries, selectQueries...)

	selectQueries = append(selectQueries,
		Limit(o.Limit), Offset(o.Offset),
		OrderBy(fmt.Sprintf("%s %s", o.Sort, o.Order)))
	for _, preload := range o.Preloads {
		selectQueries = append(selectQueries, Load(preload))
	}
	return
}

type GetArchivesResult struct {
	Archives []*modext.Archive `json:"data"`
	Total    int               `json:"total"`
	Err      error             `json:"error,omitempty"`
}

const (
	prefixGlobalAnd = "global-and"
	prefixGlobalOr  = "global-or"
)

func getArchives(or bool, opts GetArchivesOptions) (result *GetArchivesResult) {
	opts.validate()

	prefix := prefixGlobalAnd
	if or {
		prefix = prefixGlobalOr
	}

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetArchivesResult)
	}

	result = &GetArchivesResult{Archives: []*modext.Archive{}}
	defer func() {
		if len(result.Archives) > 0 || result.Total > 0 || result.Err != nil {
			Cache.RemoveWithPrefix(prefix, cacheKey)
			Cache.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	selectQueries, countQueries := opts.toQueries(or)
	archives, err := models.Archives(selectQueries...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	count, err := models.Archives(countQueries...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	result.Archives = make([]*modext.Archive, len(archives))
	result.Total = len(count)

	for i, p := range archives {
		result.Archives[i] = modext.NewArchive(p).LoadRels(p)
	}

	return
}

func GetArchives(opts GetArchivesOptions) (result *GetArchivesResult) {
	return getArchives(false, opts)
}

func SearchArchives(opts GetArchivesOptions) (result *GetArchivesResult) {
	return getArchives(true, opts)
}

func GetArchiveCount() (int64, error) {
	if c, err := Cache.Get("archive-count"); err == nil {
		return c.(int64), nil
	}

	count, err := models.Archives(Where("published_at IS NOT NULL")).CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.ErrUnknown
	}

	Cache.Set("archive-count", count, time.Hour*24*7)
	return count, nil
}

func GetArchiveStats() (size, pages uint64, err error) {
	if c, err := Cache.Get("archive-size"); err == nil {
		size = c.(uint64)
	}
	if c, err := Cache.Get("archive-pages"); err == nil {
		pages = c.(uint64)
	}

	if size > 0 && pages > 0 {
		return
	}

	archives, err := models.Archives(Where("published_at IS NOT NULL")).AllG()
	if err != nil {
		log.Println(err)
		err = errs.ErrUnknown
		return
	}

	for _, archive := range archives {
		zf, e := zip.OpenReader(archive.Path)
		if e != nil {
			log.Println(err)
			err = errs.ErrUnknown
			return
		}

		for _, f := range zf.File {
			if f.FileInfo().IsDir() {
				continue
			}
			pages++
		}
		zf.Close()

		d, e := os.Stat(archive.Path)
		if e != nil {
			log.Println(e)
			err = errs.ErrUnknown
			return
		}

		size += uint64(d.Size())
	}

	Cache.Set("archive-size", size, time.Hour*24*7)
	Cache.Set("archive-pages", pages, time.Hour*24*7)

	return
}

func PublishArchive(id int64) (*modext.Archive, error) {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrArchiveNotFound
		}
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	archive.PublishedAt = null.TimeFrom(time.Now().UTC())
	if err := archive.UpdateG(boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	// TODO: Purge cache
	return modext.NewArchive(archive), nil
}

func PublishArchives() error {
	// TODO: Purge cache
	return models.Archives(Where("published_at IS NULL")).
		UpdateAllG(models.M{"published_at": null.TimeFrom(time.Now().UTC())})
}

func UnpublishArchive(id int64) (*modext.Archive, error) {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrArchiveNotFound
		}
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	archive.PublishedAt.Valid = false
	if err := archive.UpdateG(boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	// TODO: Purge cache
	return modext.NewArchive(archive), nil
}

func UnpublishArchives() error {
	// TODO: Purge cache
	return models.Archives(Where("published_at IS NOT NULL")).
		UpdateAllG(models.M{"published_at": null.NewTime(time.Now(), false)})
}

func DeleteArchive(id int64) error {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ErrArchiveNotFound
		}
		log.Println(err)
		return errs.ErrUnknown
	}

	if err := archive.DeleteG(); err != nil {
		log.Println(err)
		return errs.ErrUnknown
	}

	// TODO: Purge cache
	os.Remove(filepath.Join(Config.Directories.Symlinks, strconv.Itoa(int(id))))
	return nil
}

func DeleteArchives() error {
	if err := models.Archives().DeleteAllG(); err != nil {
		log.Println(err)
		return errs.ErrUnknown
	}
	// TODO: Purge cache
	// TODO: Remove symlinks
	return nil
}
