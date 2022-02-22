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

var (
	ErrUnknown         = errors.New("Unknown error")
	ErrArchiveNotFound = errors.New("Archive not found")
)

func checkArchiveThumbnail(id, i, width int, w http.ResponseWriter, r *http.Request) (string, bool) {
	var fp string
	if width > 0 && width <= 1024 && width%128 == 0 {
		fp = filepath.Join(Config.Directories.Thumbnails, fmt.Sprintf("%d-%d.%d.jpg", id, i+1, width))
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
	fp, ok := checkArchiveThumbnail(id, index, width, w, r)
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

	zip, err := zip.OpenReader(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer zip.Close()

	if index > len(zip.File) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file := zip.File[index]
	d := file.FileInfo()

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

func ModerateArchives(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var artists []string
	var titles []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		arr := strings.Split(line, ":")
		if len(arr) != 2 {
			continue
		}
		if strings.EqualFold(arr[0], "artist") {
			log.Println("Blacklisting artist:", arr[1])
			artists = append(artists, arr[1])
		} else if strings.EqualFold(arr[0], "title") {
			log.Println("Blacklisting title:", arr[1])
			titles = append(titles, arr[1])
		}
	}
	f.Close()

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	archives, err := models.Archives(Load(ArchiveRels.Artists)).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	for _, archive := range archives {
		remove := false

		if archive.R != nil && len(archive.R.Artists) > 0 {
			for _, artist := range archive.R.Artists {
				if remove {
					break
				}

				for _, a := range artists {
					if strings.EqualFold(artist.Name, a) {
						artist.DeleteG()
						remove = true
						break
					}
				}
			}
		}

		if !remove {
			for _, title := range titles {
				if strings.EqualFold(archive.Title, title) {
					remove = true
					break
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
			return err
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
			return err
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
			return err
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
		return err
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
		return err
	}
	return nil
}

func CreateArchive(archive *modext.Archive) (*modext.Archive, error) {
	if archive == nil {
		return nil, nil
	} else if len(archive.Path) == 0 {
		return nil, errors.New("Archive path is required")
	}

	stat, err := os.Stat(archive.Path)
	if os.IsNotExist(err) {
		return nil, errors.New("Archive does not exist")
	}

	arc, err := models.Archives(Where("archive.path = ?", archive.Path)).OneG()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
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
		return nil, err
	}

	arc.Pages = int16(len(f.File))
	if arc.Pages > 0 {
		d := f.File[0].FileInfo()
		arc.CreatedAt = d.ModTime()
		arc.UpdatedAt = arc.CreatedAt
	}
	f.Close()

	if err := arc.InsertG(boil.Infer()); err != nil {
		return nil, err
	} else if err := refreshArchiveRels(arc, archive); err != nil {
		return nil, err
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
		result.Err = ErrUnknown
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
	selectQueries = append(selectQueries, GroupBy("archive.id"))
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

	cacheKey := fmt.Sprintf("%s%v", makeCacheKey(opts))
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
		result.Err = ErrUnknown
		return
	}

	count, err := models.Archives(countQueries...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
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

func PublishArchive(id int64) (*modext.Archive, error) {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrArchiveNotFound
		}
		return nil, ErrUnknown
	}

	archive.PublishedAt = null.TimeFrom(time.Now().UTC())
	if err := archive.UpdateG(boil.Infer()); err != nil {
		return nil, ErrUnknown
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
			return nil, ErrArchiveNotFound
		}
		return nil, ErrUnknown
	}

	archive.PublishedAt.Valid = false
	if err := archive.UpdateG(boil.Infer()); err != nil {
		return nil, ErrUnknown
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
			return ErrArchiveNotFound
		}
		return ErrUnknown
	}

	if err := archive.DeleteG(); err != nil {
		return ErrUnknown
	}

	// TODO: Purge cache
	os.Remove(filepath.Join(Config.Directories.Symlinks, strconv.Itoa(int(id))))
	return nil
}

func DeleteArchives() error {
	// TODO: Purge cache
	// TODO: Remove symlinks
	return models.Archives().DeleteAllG()
}
