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

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CreateCircle(name string) (*modext.Circle, error) {
	name = strings.Title(strings.TrimSpace(name))
	if len(name) == 0 {
		return nil, errs.ErrCircleNameRequired
	} else if len(name) > 128 {
		return nil, errs.ErrCircleNameTooLong
	}

	slug := slugify(name)
	circle, err := models.Circles(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		circle = &models.Circle{
			Name: name,
			Slug: slug,
		}
		if err = circle.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.ErrUnknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	return modext.NewCircle(circle), nil
}

func GetCircle(slug string) (*modext.Circle, error) {
	circle, err := models.Circles(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrCircleNotFound
		}
		log.Println(err)
		return nil, errs.ErrUnknown
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
		Select("circle.*", "COUNT(archive.circle_id) AS archive_count"),
		InnerJoin("archive_circles archive ON archive.circle_id = circle.id"),
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
		result.Err = errs.ErrUnknown
		return
	}

	count, err := models.Circles().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	result.Total = int(count)
	return
}

func GetCircleCount() (int64, error) {
	if c, err := Cache.Get("circle-count"); err == nil {
		return c.(int64), nil
	}

	count, err := models.Circles().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.ErrUnknown
	}

	Cache.Set("circle-count", count, time.Hour*24*7)
	return count, nil
}

var isCircleValidMap = QueryMapCache{
	Map: make(map[string]bool),
}

func IsCircleValid(str string) (isValid bool) {
	str = slugify(str)

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
