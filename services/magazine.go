package services

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	. "koushoku/cache"
	"koushoku/errs"

	"koushoku/models"
	"koushoku/modext"

	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CreateMagazine(name string) (*modext.Magazine, error) {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return nil, errs.ErrMagazineNameRequired
	} else if len(name) > 128 {
		return nil, errs.ErrMagazineNameTooLong
	}

	slug := slug.Make(name)

	magazine, err := models.Magazines(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		magazine = &models.Magazine{
			Name: name,
			Slug: slug,
		}
		if err = magazine.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.ErrUnknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	return modext.NewMagazine(magazine), nil
}

func GetMagazine(slug string) (*modext.Magazine, error) {
	magazine, err := models.Magazines(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrMagazineNotFound
		}
		log.Println(err)
		return nil, errs.ErrUnknown
	}
	return modext.NewMagazine(magazine), nil
}

type GetMagazinesOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetMagazinesResult struct {
	Magazines []*modext.Magazine
	Total     int
	Err       error
}

const prefixgm = "magazines"

func GetMagazines(opts GetMagazinesOptions) (result *GetMagazinesResult) {
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(prefixgm, cacheKey); err == nil {
		return c.(*GetMagazinesResult)
	}

	result = &GetMagazinesResult{Magazines: []*modext.Magazine{}}
	defer func() {
		if len(result.Magazines) > 0 || result.Total > 0 || result.Err != nil {
			Cache.RemoveWithPrefix(prefixgm, cacheKey)
			Cache.SetWithPrefix(prefixgm, cacheKey, result, time.Hour*24*7)
		}
	}()

	q := []QueryMod{
		Select("magazine.*", "COUNT(archive.id) AS archive_count"),
		InnerJoin("archive ON archive.magazine_id = magazine.id"),
		GroupBy("magazine.id"),
		OrderBy("magazine.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Magazines(q...).BindG(context.Background(), &result.Magazines)
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	count, err := models.Magazines().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
	}

	result.Total = int(count)
	return
}

var isMagazineValidMap = QueryMapCache{
	Map: make(map[string]bool),
}

func IsMagazineValid(str string) (isValid bool) {
	str = slug.Make(str)

	isMagazineValidMap.RLock()
	v, ok := isMagazineValidMap.Map[str]
	isMagazineValidMap.RUnlock()

	if ok {
		return v
	}

	result := GetMagazines(GetMagazinesOptions{})
	if result.Err != nil {
		return
	}

	defer func() {
		isMagazineValidMap.Lock()
		isMagazineValidMap.Map[str] = isValid
		isMagazineValidMap.Unlock()
	}()

	for _, magazine := range result.Magazines {
		if magazine.Slug == str {
			isValid = true
			break
		}
	}
	return
}
