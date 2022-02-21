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

func CreateParody(name string) (*modext.Parody, error) {
	name = strings.Title(strings.TrimSpace(name))

	if len(name) == 0 {
		return nil, errors.New("Parody name is required")
	} else if len(name) > 128 {
		return nil, errors.New("Parody name is too long")
	}

	slug := slug.Make(name)

	parody, err := models.Parodies(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		parody = &models.Parody{
			Name: name,
			Slug: slug,
		}
		if err = parody.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, err
		}
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return modext.NewParody(parody), nil
}

func GetParody(slug string) (*modext.Parody, error) {
	parody, err := models.Parodies(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Parody does not exist")
		}
		return nil, err
	}
	return modext.NewParody(parody), nil
}

type GetParodiesOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetParodiesResult struct {
	Parodies []*modext.Parody
	Total    int
	Err      error
}

const prefixgp = "parodies"

func GetParodies(opts GetParodiesOptions) (result *GetParodiesResult) {
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(prefixgp, cacheKey); err == nil {
		return c.(*GetParodiesResult)
	}

	result = &GetParodiesResult{Parodies: []*modext.Parody{}}
	defer func() {
		if len(result.Parodies) > 0 || result.Total > 0 || result.Err != nil {
			Cache.RemoveWithPrefix(prefixgp, cacheKey)
			Cache.SetWithPrefix(prefixgp, cacheKey, result, time.Hour*24*7)
		}
	}()

	q := []QueryMod{
		Select("parody.*", "COUNT(archive.id) AS archive_count"),
		InnerJoin("archive ON archive.parody_id = parody.id"),
		GroupBy("parody.id"),
		OrderBy("parody.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Parodies(q...).BindG(context.Background(), &result.Parodies)
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
		return
	}

	count, err := models.Parodies().CountG()
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
		return
	}

	result.Total = int(count)
	return
}

var isParodyValidMap = QueryMapCache{
	Map: make(map[string]bool),
}

func IsParodyValid(str string) (isValid bool) {
	str = slug.Make(str)

	isParodyValidMap.RLock()
	v, ok := isParodyValidMap.Map[str]
	isParodyValidMap.RUnlock()

	if ok {
		return v
	}

	result := GetParodies(GetParodiesOptions{})
	if result.Err != nil {
		return
	}

	defer func() {
		isParodyValidMap.Lock()
		isParodyValidMap.Map[str] = isValid
		isParodyValidMap.Unlock()
	}()

	for _, parody := range result.Parodies {
		if parody.Slug == str {
			isValid = true
			break
		}
	}
	return
}
