package services

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	. "koushoku/cache"
	. "koushoku/config"

	"koushoku/database"
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

func init() {
	slug.CustomSub = make(map[string]string)
	slug.CustomSub["☆"] = "-"
}

var archiveRgx = regexp.MustCompile(`(\(|\[|\{)?[^\(\[\{\}\]\)]+(\}\)|\])?`)
var replRgx = regexp.MustCompile(`(?i)(fakku|irodori comics|x?\d+00x?)`)

func populateArchive(archive *modext.Archive) error {
	if archive == nil {
		return nil
	}

	fileName := FileName(archive.Path)
	if stat, err := os.Stat(archive.Path); err == nil {
		archive.Size = stat.Size()
	} else {
		return err
	}

	var (
		artists   = make(map[string]string)
		circles   = make(map[string]string)
		magazines = make(map[string]string)
		parodies  = make(map[string]string)
		tags      = make(map[string]string)
	)

	if metadata, ok := metadataMap.Map[slugify(fileName)]; ok {
		for _, parody := range metadata.Parodies {
			slug := slugify(parody)
			if v, ok := alias.Parodies[slug]; ok {
				slug = slugify(v)
				parody = v
			}
			parodies[slug] = parody
		}

		for _, tag := range metadata.Tags {
			slug := slugify(tag)
			if v, ok := alias.Tags[slug]; ok {
				tag = v
				slug = slugify(v)
			}
			if _, ok := blacklist.Tags[slug]; ok {
				return nil
			}
			tags[slug] = tag
		}
	}

	matches := archiveRgx.FindAllString(fileName, -1)
	if len(matches) == 0 {
		return nil
	}

	var title string
	for i, match := range matches {
		match = strings.TrimSpace(match)
		if len(match) == 0 {
			continue
		}

		if strings.HasPrefix(match, "[") {
			if i == 0 {
				match = strings.TrimSuffix(strings.TrimPrefix(match, "["), "]")
				if match = strings.TrimSpace(match); len(match) == 0 {
					continue
				}

				names := strings.Split(match, ",")
				for _, name := range names {
					if name = strings.TrimSpace(name); len(name) > 0 {
						artists[slugify(name)] = name
					}
				}
			}

		} else if strings.HasPrefix(match, "(") {
			if i == 1 {
				match = strings.TrimSuffix(strings.TrimPrefix(match, "("), ")")
				if match = strings.TrimSpace(match); len(match) == 0 {
					continue
				}

				if len(artists) > 0 {
					for k, v := range artists {
						circles[k] = v
						delete(artists, k)
					}
				}

				names := strings.Split(match, ",")
				for _, name := range names {
					if name = strings.TrimSpace(name); len(name) > 0 {
						artists[slugify(name)] = name
					}
				}
			} else if i == 2 || i == 3 {
				match = strings.TrimSuffix(strings.TrimPrefix(match, "("), ")")
				if match = strings.TrimSpace(match); len(match) == 0 {
					continue
				}

				if i < len(matches)-1 {
					next := matches[i+1]
					if len(next) > 0 &&
						!(strings.HasPrefix(match, "[") ||
							strings.HasPrefix(match, "(") ||
							strings.HasPrefix(next, "{")) {
						continue
					}
				}

				if strings.HasPrefix(match, "x") ||
					strings.EqualFold(match, "temp") ||
					strings.EqualFold(match, "strong") ||
					strings.EqualFold(match, "complete") {
					continue
				}

				names := strings.Split(match, ", ")
				for _, name := range names {
					if name = strings.TrimSpace(name); len(name) > 0 {
						magazines[slugify(name)] = name
					}
				}
			}
		} else if strings.HasPrefix(match, "{") {
			match = strings.TrimSuffix(strings.TrimPrefix(match, "{"), "}")
			match = strings.TrimSpace(match)

			if len(match) == 0 ||
				strings.Contains(slugify(match), "comic") ||
				strings.Contains(slugify(match), "2d-market") {
				continue
			}

			match = strings.ReplaceAll(match, "zero gravity", "zero-gravity")
			match = strings.ReplaceAll(match, "dark skin", "dark-skin")
			match = strings.ReplaceAll(match, "heart pupil", "heart-pupil")

			names := strings.Split(match, " ")
			for _, name := range names {
				name = strings.TrimSpace(name)
				if len(name) > 0 {
					tags[slugify(name)] = strings.ReplaceAll(name, "-", " ")
				}
			}
		} else if i == 1 || i == 2 {
			match = strings.TrimSpace(replRgx.ReplaceAllString(match, ""))
			if len(match) > 0 {
				title = match
			}
		}
	}

	titleSlug := slugify(title)
	if v, ok := alias.Archives[titleSlug]; ok {
		titleSlug = slugify(title)
		title = v
	}

	if _, ok := blacklist.Archives[titleSlug]; ok {
		return nil
	}

	for _, v := range blacklist.ArchivesG {
		if strings.Contains(titleSlug, v) {
			return nil
		}
	}

	for slug, artist := range artists {
		if v, ok := alias.Artists[slug]; ok {
			slug = slugify(v)
			artist = v
		}
		if _, ok := blacklist.Artists[slug]; ok {
			return nil
		}
		archive.Artists = append(archive.Artists,
			&modext.Artist{Slug: slug, Name: artist})
	}

	for slug, circle := range circles {
		if v, ok := alias.Circles[slug]; ok {
			slug = slugify(v)
			circle = v
		}
		if _, ok := blacklist.Circles[slug]; ok {
			return nil
		}
		archive.Circles = append(archive.Circles,
			&modext.Circle{Slug: slug, Name: circle})
	}

	for slug, magazine := range magazines {
		if v, ok := alias.Magazines[slug]; ok {
			slug = slugify(v)
			magazine = v
		}
		if _, ok := blacklist.Magazines[slug]; ok {
			return nil
		}
		archive.Magazines = append(archive.Magazines,
			&modext.Magazine{Slug: slug, Name: magazine})
	}

	for slug, parody := range parodies {
		if v, ok := alias.Parodies[slug]; ok {
			slug = slugify(v)
			parody = v
		}
		archive.Parodies = append(archive.Parodies,
			&modext.Parody{Slug: slug, Name: parody})
	}

	for slug, tag := range tags {
		if v, ok := alias.Tags[slug]; ok {
			slug = slugify(v)
			tag = v
		}
		if _, ok := blacklist.Tags[slug]; ok {
			return nil
		}

		isDuplicate := false
		for _, t := range archive.Tags {
			if slug == slugify(t.Name) {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			archive.Tags = append(archive.Tags,
				&modext.Tag{Slug: slug, Name: tag})
		}
	}

	zf, err := zip.OpenReader(archive.Path)
	if err != nil {
		if err == zip.ErrFormat {
			log.Println(err, archive.Path)
			return nil
		}
		return err
	}
	defer zf.Close()

	for _, f := range zf.File {
		stat := f.FileInfo()
		name := stat.Name()

		if stat.IsDir() ||
			!(strings.HasSuffix(name, ".jpg") ||
				strings.HasSuffix(name, ".png")) {
			continue
		}

		archive.Pages++
		if archive.CreatedAt == 0 {
			archive.CreatedAt = stat.ModTime().Unix()
		}
	}

	if archive.Pages < 3 {
		return nil
	}

	archive.Title = title
	archive.Slug = titleSlug

	return nil
}

var pathBlacklist = []string{
	"/cover",
	"/doujin",
	"/illustration",
	"/interview",
	"/non-h",
	"/spread",
	"/western",
	"/D/",
}

func isPathBlacklisted(path string) bool {
	for _, p := range pathBlacklist {
		if strings.Contains(path, p) {
			return true
		}
	}
	return !(strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".cbz"))
}

func IndexArchives() {
	if _, err := os.Stat(Config.Directories.Symlinks); os.IsNotExist(err) {
		if err := os.MkdirAll(Config.Directories.Symlinks, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	initAlias()
	initBlacklist()
	initMetadata()

	var files []string
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || isPathBlacklisted(path) {
			return err
		}

		log.Println("Found archive", filepath.Base(path))
		files = append(files, path)
		return nil
	}

	if err := filepath.Walk(Config.Directories.Data, walkFn); err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	c := make(chan bool, 20)
	defer func() {
		close(c)
	}()

	var archives []*modext.Archive
	var mu sync.Mutex

	for _, path := range files {
		c <- true

		go func(path string) {
			defer func() {
				wg.Done()
				<-c
			}()

			archive := &modext.Archive{Path: path}
			log.Println("Populating archive", filepath.Base(path))
			if err := populateArchive(archive); err != nil {
				log.Fatalln(err)
			}

			if len(archive.Title) > 0 {
				mu.Lock()
				archives = append(archives, archive)
				mu.Unlock()
			}
		}(path)
	}
	wg.Wait()

	sort.SliceStable(archives, func(i, j int) bool {
		return archives[i].CreatedAt < archives[j].CreatedAt
	})

	wg.Add(len(archives))
	for _, archive := range archives {
		c <- true
		go func(archive *modext.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			log.Println("Indexing archive", filepath.Base(archive.Path))
			c, err := CreateArchive(archive)
			if c != nil && err == nil {
				CreateArchiveSymlink(c)
			}
		}(archive)
	}
	wg.Wait()
}

func ModerateArchives() {
	initBlacklist()

	archives, err := models.Archives(
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	for _, archive := range archives {
		titleSlug := slugify(archive.Title)
		_, remove := blacklist.Archives[titleSlug]

		if archive.R != nil && len(archive.R.Artists) > 0 {
			for _, artist := range archive.R.Artists {
				if _, ok := blacklist.Artists[slugify(artist.Name)]; ok {
					artist.DeleteG()
					remove = true
				}
			}
		}

		if !remove {
			for _, t := range blacklist.ArchivesG {
				if strings.Contains(titleSlug, t) {
					remove = true
					break
				}
			}
		}

		if !remove && archive.R != nil && len(archive.R.Tags) > 0 {
			for _, tag := range archive.R.Tags {
				if _, ok := blacklist.Tags[slugify(tag.Name)]; ok {
					tag.DeleteG()
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

var refreshArchiveRelsCache struct {
	Artists   map[string]*models.Artist
	Circles   map[string]*models.Circle
	Magazines map[string]*models.Magazine
	Parodies  map[string]*models.Parody
	Tags      map[string]*models.Tag

	sync.RWMutex
	sync.Once
}

func refreshArchiveRels(e boil.Executor, arc *models.Archive, archive *modext.Archive) error {
	var err error
	refreshArchiveRelsCache.Do(func() {
		refreshArchiveRelsCache.Lock()
		defer refreshArchiveRelsCache.Unlock()

		refreshArchiveRelsCache.Artists = make(map[string]*models.Artist)
		refreshArchiveRelsCache.Circles = make(map[string]*models.Circle)
		refreshArchiveRelsCache.Magazines = make(map[string]*models.Magazine)
		refreshArchiveRelsCache.Parodies = make(map[string]*models.Parody)
		refreshArchiveRelsCache.Tags = make(map[string]*models.Tag)
	})

	if len(archive.Artists) > 0 {
		var artists []*models.Artist
		for _, artist := range archive.Artists {
			refreshArchiveRelsCache.RLock()
			model, ok := refreshArchiveRelsCache.Artists[artist.Name]
			refreshArchiveRelsCache.RUnlock()

			if ok {
				artists = append(artists, model)
				continue
			}

			refreshArchiveRelsCache.Lock()
			artist, err = CreateArtist(artist.Name)
			if err != nil {
				refreshArchiveRelsCache.Unlock()
				return err
			}

			model = &models.Artist{
				ID:   artist.ID,
				Slug: artist.Slug,
				Name: artist.Name,
			}
			refreshArchiveRelsCache.Artists[artist.Name] = model
			refreshArchiveRelsCache.Unlock()

			artists = append(artists, model)
		}
		if err := arc.SetArtists(e, false, artists...); err != nil {
			log.Println(err)
			return errs.ErrUnknown
		}
	}

	if len(archive.Circles) > 0 {
		var circles []*models.Circle
		for _, circle := range archive.Circles {
			refreshArchiveRelsCache.RLock()
			model, ok := refreshArchiveRelsCache.Circles[circle.Name]
			refreshArchiveRelsCache.RUnlock()

			if ok {
				circles = append(circles, model)
				continue
			}

			refreshArchiveRelsCache.Lock()
			circle, err := CreateCircle(circle.Name)
			if err != nil {
				refreshArchiveRelsCache.Unlock()
				return err
			}

			model = &models.Circle{
				ID:   circle.ID,
				Slug: circle.Slug,
				Name: circle.Name,
			}
			refreshArchiveRelsCache.Circles[circle.Name] = model
			refreshArchiveRelsCache.Unlock()

			circles = append(circles, model)
		}
		if err := arc.SetCircles(e, false, circles...); err != nil {
			var str []string
			for _, circle := range circles {
				str = append(str, circle.Name)
			}
			log.Println(err, str)
			return errs.ErrUnknown
		}
	}

	if len(archive.Magazines) > 0 {
		var magazines []*models.Magazine
		for _, magazine := range archive.Magazines {
			refreshArchiveRelsCache.RLock()
			model, ok := refreshArchiveRelsCache.Magazines[magazine.Name]
			refreshArchiveRelsCache.RUnlock()

			if ok {
				magazines = append(magazines, model)
				continue
			}

			refreshArchiveRelsCache.Lock()
			magazine, err := CreateMagazine(magazine.Name)
			if err != nil {
				refreshArchiveRelsCache.Unlock()
				return err
			}

			model = &models.Magazine{
				ID:   magazine.ID,
				Slug: magazine.Slug,
				Name: magazine.Name,
			}
			refreshArchiveRelsCache.Magazines[magazine.Name] = model
			refreshArchiveRelsCache.Unlock()

			magazines = append(magazines, model)
		}
		if err := arc.SetMagazines(e, false, magazines...); err != nil {
			var str []string
			for _, magazine := range magazines {
				str = append(str, magazine.Name)
			}
			log.Println(err, str)
			return errs.ErrUnknown
		}
	}

	if len(archive.Parodies) > 0 {
		var parodies []*models.Parody
		for _, parody := range archive.Parodies {
			refreshArchiveRelsCache.RLock()
			model, ok := refreshArchiveRelsCache.Parodies[parody.Name]
			refreshArchiveRelsCache.RUnlock()

			if ok {
				parodies = append(parodies, model)
				continue
			}

			refreshArchiveRelsCache.Lock()
			parody, err := CreateParody(parody.Name)
			if err != nil {
				refreshArchiveRelsCache.Unlock()
				return err
			}

			model = &models.Parody{
				ID:   parody.ID,
				Slug: parody.Slug,
				Name: parody.Name,
			}
			refreshArchiveRelsCache.Parodies[parody.Name] = model
			refreshArchiveRelsCache.Unlock()

			parodies = append(parodies, model)
		}
		if err := arc.SetParodies(e, false, parodies...); err != nil {
			var str []string
			for _, parody := range parodies {
				str = append(str, parody.Name)
			}
			log.Println(err, str)
			return errs.ErrUnknown
		}
	}

	if len(archive.Tags) > 0 {
		var tags []*models.Tag
		for _, tag := range archive.Tags {
			refreshArchiveRelsCache.RLock()
			model, ok := refreshArchiveRelsCache.Tags[tag.Name]
			refreshArchiveRelsCache.RUnlock()

			if ok {
				tags = append(tags, model)
				continue
			}

			refreshArchiveRelsCache.Lock()
			tag, err := CreateTag(tag.Name)
			if err != nil {
				refreshArchiveRelsCache.Unlock()
				return err
			}

			model = &models.Tag{
				ID:   tag.ID,
				Slug: tag.Slug,
				Name: tag.Name,
			}
			refreshArchiveRelsCache.Tags[tag.Name] = model
			refreshArchiveRelsCache.Unlock()

			tags = append(tags, model)
		}
		if err := arc.SetTags(e, false, tags...); err != nil {
			var str []string
			for _, tag := range tags {
				str = append(str, tag.Name)
			}
			log.Println(err, str)
			return errs.ErrUnknown
		}
	}

	return nil
}

func CreateArchive(archive *modext.Archive) (*modext.Archive, error) {
	if archive == nil {
		return nil, nil
	} else if len(archive.Path) == 0 {
		return nil, errs.ErrArchivePathRequired
	}

	selectQueries := []QueryMod{
		Where("archive.slug = ?", archive.Slug),
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Circles),
		Load(ArchiveRels.Magazines),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	}
	if len(archive.Artists) > 0 {
		var artists []string
		for _, artist := range archive.Artists {
			artists = append(artists, slugify(artist.Name))
		}

		mods, query, args := joinRelQuery(true, models.TableNames.Artist, artists)
		selectQueries = append(selectQueries, mods...)
		selectQueries = append(selectQueries, Where(query, args...))
	} else if len(archive.Magazines) > 0 {
		var magazines []string
		for _, magazine := range archive.Magazines {
			magazines = append(magazines, slugify(magazine.Name))
		}

		mods, query, args := joinRelQuery(true, models.TableNames.Magazine, magazines)
		selectQueries = append(selectQueries, mods...)
		selectQueries = append(selectQueries, Where(query, args...))
	} else if len(archive.Circles) > 0 {
		var circles []string
		for _, circle := range archive.Circles {
			circles = append(circles, slugify(circle.Name))
		}

		mods, query, args := joinRelQuery(true, models.TableNames.Circle, circles)
		selectQueries = append(selectQueries, mods...)
		selectQueries = append(selectQueries, Where(query, args...))
	} else {
		selectQueries = append(selectQueries,
			Where("archive.path = ?", archive.Path))
	}

	arc, err := models.Archives(selectQueries...).OneG()
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	tx, err := database.Conn.Begin()
	if err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	var isDuplicate bool
	if arc == nil {
		arc = &models.Archive{
			Title: archive.Title,
			Slug:  archive.Slug,
		}
		if archive.CreatedAt > 0 {
			arc.CreatedAt = time.Unix(archive.CreatedAt, 0)
			arc.UpdatedAt = arc.CreatedAt
		}
	} else {
		isDuplicate = true
		arc.UpdatedAt = time.Unix(archive.CreatedAt, 0)
	}

	arc.Path = archive.Path
	arc.Pages = archive.Pages
	arc.Size = archive.Size

	upsert := arc.Insert
	if isDuplicate {
		upsert = arc.Update
	}

	if err := upsert(tx, boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	} else if err := refreshArchiveRels(tx, arc, archive); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	// TODO: Purge cache
	return modext.NewArchive(arc), nil
}

func CreateArchiveSymlink(archive *modext.Archive) error {
	if archive == nil {
		return nil
	}

	symlink := filepath.Join(Config.Directories.Symlinks, strconv.Itoa(int(archive.ID)))
	return os.Symlink(archive.Path, symlink)
}

func validateRels(rels []string) (result []string) {
	for _, v := range rels {
		if strings.EqualFold(v, ArchiveRels.Artists) {
			result = append(result, ArchiveRels.Artists)
		} else if strings.EqualFold(v, ArchiveRels.Circles) {
			result = append(result, ArchiveRels.Circles)
		} else if strings.EqualFold(v, ArchiveRels.Magazines) {
			result = append(result, ArchiveRels.Magazines)
		} else if strings.EqualFold(v, ArchiveRels.Parodies) {
			result = append(result, ArchiveRels.Parodies)
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
	Path  string `json:"0,omitempty"`
	Title string `json:"1,omitempty"`

	Circles   []string `json:"2,omitempty"`
	Magazines []string `json:"3,omitempty"`
	Parodies  []string `json:"4,omitempty"`
	Artists   []string `json:"5,omitempty"`
	Tags      []string `json:"6,omitempty"`

	Limit    int      `json:"7,omitempty"`
	Offset   int      `json:"8,omitempty"`
	Preloads []string `json:"9,omitempty"`
	Sort     string   `json:"10,omitempty"`
	Order    string   `json:"11,omitempty"`
}

func (o *GetArchivesOptions) validate() {
	o.Path = strings.ToLower(o.Path)
	o.Title = slugify(o.Title)

	for i, circle := range o.Circles {
		o.Circles[i] = slugify(circle)
	}
	sort.Strings(o.Circles)

	for i, magazine := range o.Magazines {
		o.Magazines[i] = slugify(magazine)
	}
	sort.Strings(o.Magazines)

	for i, parody := range o.Parodies {
		o.Parodies[i] = slugify(parody)
	}
	sort.Strings(o.Parodies)

	for i, artist := range o.Artists {
		o.Artists[i] = slugify(artist)
	}
	sort.Strings(o.Artists)

	for i, tag := range o.Tags {
		o.Tags[i] = slugify(tag)
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

func joinRelQuery(isOr bool, n string, values []string) ([]QueryMod, string, []interface{}) {
	var queries []QueryMod
	var query string
	var args []interface{}

	queries = append(queries,
		InnerJoin(fmt.Sprintf("archive_%s j ON j.archive_id = archive.id", pluralize(n))),
		InnerJoin(fmt.Sprintf("%s ON %[1]s.id = j.%[1]s_id", n)))

	var q []string
	for _, v := range values {
		if isOr {
			q = append(q, n+".slug ILIKE '%' || ? || '%'")
		} else {
			q = append(q, n+".slug = ?")
		}
		args = append(args, v)
	}

	if len(q) > 1 {
		query = fmt.Sprintf("(%s)", strings.Join(q, " OR "))
	} else {
		query = q[0]
	}
	return queries, query, args
}

func (o *GetArchivesOptions) toQueries(isOr bool) (selectQueries, countQueries []QueryMod) {
	countQueries = []QueryMod{Select("1")}

	var queries []string
	var arguments []interface{}

	if len(o.Path) > 0 {
		queries = append(queries, "archive.path ILIKE '%' || ? || '%'")
		arguments = append(arguments, o.Path)
	}

	if len(o.Title) > 0 {
		queries = append(queries, "archive.slug ILIKE '%' || ? || '%'")
		arguments = append(arguments, o.Title)
	}

	if len(o.Artists) > 0 {
		mods, query, args := joinRelQuery(isOr, models.TableNames.Artist, o.Artists)
		selectQueries = append(selectQueries, mods...)
		queries = append(queries, query)
		arguments = append(arguments, args...)
	}

	if len(o.Circles) > 0 {
		mods, query, args := joinRelQuery(isOr, models.TableNames.Circle, o.Circles)
		selectQueries = append(selectQueries, mods...)
		queries = append(queries, query)
		arguments = append(arguments, args...)
	}

	if len(o.Magazines) > 0 {
		mods, query, args := joinRelQuery(isOr, models.TableNames.Magazine, o.Magazines)
		selectQueries = append(selectQueries, mods...)
		queries = append(queries, query)
		arguments = append(arguments, args...)
	}

	if len(o.Parodies) > 0 {
		mods, query, args := joinRelQuery(isOr, models.TableNames.Parody, o.Parodies)
		selectQueries = append(selectQueries, mods...)
		queries = append(queries, query)
		arguments = append(arguments, args...)
	}

	if len(o.Tags) > 0 {
		mods, query, args := joinRelQuery(false, models.TableNames.Tag, o.Tags)
		selectQueries = append(selectQueries, mods...)
		queries = append(queries, query)
		arguments = append(arguments, args...)
	}

	if len(o.Artists) > 0 ||
		len(o.Circles) > 0 ||
		len(o.Magazines) > 0 ||
		len(o.Parodies) > 0 ||
		len(o.Tags) > 0 {
		selectQueries = append(selectQueries, GroupBy("archive.id"))
	}

	if len(queries) > 0 {
		op := " AND "
		if isOr {
			op = " OR "
		}
		selectQueries = append(selectQueries,
			Where(strings.Join(queries, op), arguments...))
	}

	selectQueries = append(selectQueries,
		Where("archive.published_at IS NOT NULL"))
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

func GetArchiveStats() (size, pages int64, err error) {
	if c, err := Cache.Get("archive-size"); err == nil {
		size = c.(int64)
	}
	if c, err := Cache.Get("archive-pages"); err == nil {
		pages = c.(int64)
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
		pages += int64(archive.Pages)
		size += archive.Size
	}

	Cache.Set("archive-size", size, time.Hour*24*7)
	Cache.Set("archive-pages", pages, time.Hour*24*7)

	return
}

func GetArchiveSymlink(id int) (string, error) {
	symlink := filepath.Join(Config.Directories.Symlinks, strconv.Itoa(id))
	return os.Readlink(symlink)
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

func PurgeArchiveThumbnails() {
	if err := os.RemoveAll(Config.Directories.Thumbnails); err != nil {
		log.Fatalln(err)
	} else if err := os.MkdirAll(Config.Directories.Thumbnails, 0755); err != nil {
		log.Fatalln(err)
	}
}

func PurgeArchiveSymlinks() {
	if err := os.RemoveAll(Config.Directories.Symlinks); err != nil {
		log.Fatalln(err)
	}
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
