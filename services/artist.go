package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	. "koushoku/cache"

	"koushoku/models"
	"koushoku/modext"

	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"

	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CreateArtist(name string) (*modext.Artist, error) {
	name = strings.Title(strings.TrimSpace(name))

	if len(name) == 0 {
		return nil, errors.New("Artist name is required")
	} else if len(name) > 128 {
		return nil, errors.New("Artist name must be at most 128 characters")
	}

	slug := slug.Make(name)

	artist, err := models.Artists(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		artist = &models.Artist{
			Name: name,
			Slug: slug,
		}
		if err = artist.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, err
		}
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return modext.NewArtist(artist), nil
}

func GetArtist(slug string) (*modext.Artist, error) {
	artist, err := models.Artists(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Artist does not exist")
		}
		return nil, err
	}
	return modext.NewArtist(artist), nil
}

type GetArtistsOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetArtistsResult struct {
	Artists []*modext.Artist
	Total   int
	Err     error
}

const prefixgart = "artists"

func GetArtists(opts GetArtistsOptions) (result *GetArtistsResult) {
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(prefixgart, cacheKey); err == nil {
		return c.(*GetArtistsResult)
	}

	result = &GetArtistsResult{Artists: []*modext.Artist{}}
	defer func() {
		if len(result.Artists) > 0 || result.Total > 0 || result.Err != nil {
			Cache.RemoveWithPrefix(prefixgart, cacheKey)
			Cache.SetWithPrefix(prefixgart, cacheKey, result, time.Hour*24*7)
		}
	}()

	q := []QueryMod{
		Select("artist.*", "COUNT(archive.id) AS archive_count"),
		InnerJoin("archive_artists ar ON ar.artist_id = artist.id"),
		InnerJoin("archive ON archive.id = ar.archive_id"),
		GroupBy("artist.id"),
		OrderBy("artist.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Artists(q...).BindG(context.Background(), &result.Artists)
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
		return
	}

	count, err := models.Artists().CountG()
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
		return
	}

	result.Total = int(count)
	return
}

var isArtistValidMap = QueryMapCache{
	Map: make(map[string]bool),
}

func IsArtistValid(str string) (isValid bool) {
	str = slug.Make(str)

	isArtistValidMap.RLock()
	v, ok := isArtistValidMap.Map[str]
	isArtistValidMap.RUnlock()

	if ok {
		return v
	}

	result := GetArtists(GetArtistsOptions{})
	if result.Err != nil {
		return
	}

	defer func() {
		isArtistValidMap.Lock()
		isArtistValidMap.Map[str] = isValid
		isArtistValidMap.Unlock()
	}()

	for _, artist := range result.Artists {
		if artist.Slug == str {
			isValid = true
			break
		}
	}
	return
}
