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

func CreateCircle(name string) (*modext.Circle, error) {
	name = strings.Title(strings.TrimSpace(name))

	if len(name) == 0 {
		return nil, errors.New("Circle name is required")
	} else if len(name) > 128 {
		return nil, errors.New("Circle name is too long")
	}

	slug := slug.Make(name)

	circle, err := models.Circles(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		circle = &models.Circle{
			Name: name,
			Slug: slug,
		}
		if err = circle.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, err
		}
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return modext.NewCircle(circle), nil
}

func GetCircle(slug string) (*modext.Circle, error) {
	circle, err := models.Circles(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Circle does not exist")
		}
		return nil, err
	}
	return modext.NewCircle(circle), nil
}

type GetCirclesOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetCirclesResult struct {
	Circles []*modext.Circle
	Total   int
	Err     error
}

const prefixgc = "circles"

func GetCircles(opts GetCirclesOptions) (result *GetCirclesResult) {
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(prefixgc, cacheKey); err == nil {
		return c.(*GetCirclesResult)
	}

	result = &GetCirclesResult{Circles: []*modext.Circle{}}
	defer func() {
		if len(result.Circles) > 0 || result.Total > 0 || result.Err != nil {
			Cache.RemoveWithPrefix(prefixgc, cacheKey)
			Cache.SetWithPrefix(prefixgc, cacheKey, result, time.Hour*24*7)
		}
	}()

	q := []QueryMod{
		Select("circle.*", "COUNT(archive.id) AS archive_count"),
		InnerJoin("archive ON archive.circle_id = circle.id"),
		GroupBy("circle.id"),
		OrderBy("circle.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Circles(q...).BindG(context.Background(), &result.Circles)
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
		return
	}

	count, err := models.Circles().CountG()
	if err != nil {
		log.Println(err)
		result.Err = ErrUnknown
		return
	}

	result.Total = int(count)
	return
}

var isCircleValidMap = QueryMapCache{
	Map: make(map[string]bool),
}

func IsCircleValid(str string) (isValid bool) {
	str = slug.Make(str)

	isCircleValidMap.RLock()
	v, ok := isCircleValidMap.Map[str]
	isCircleValidMap.RUnlock()

	if ok {
		return v
	}

	result := GetCircles(GetCirclesOptions{})
	if result.Err != nil {
		return
	}

	defer func() {
		isCircleValidMap.Lock()
		isCircleValidMap.Map[str] = isValid
		isCircleValidMap.Unlock()
	}()

	for _, circle := range result.Circles {
		if circle.Slug == str {
			isValid = true
			break
		}
	}
	return
}
