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

func CreateTag(name string) (*modext.Tag, error) {
	name = strings.Title(strings.TrimSpace(name))
	if len(name) == 0 {
		return nil, errs.ErrTagNameRequired
	} else if len(name) > 128 {
		return nil, errs.ErrTagNameTooLong
	}

	slug := slugify(name)

	tag, err := models.Tags(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		tag = &models.Tag{
			Name: name,
			Slug: slug,
		}
		if err = tag.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.ErrUnknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.ErrUnknown
	}

	return modext.NewTag(tag), nil
}

func GetTag(slug string) (*modext.Tag, error) {
	tag, err := models.Tags(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrTagNotFound
		}
		log.Println(err)
		return nil, errs.ErrUnknown
	}
	return modext.NewTag(tag), nil
}

type GetTagsOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetTagsResult struct {
	Tags  []*modext.Tag
	Total int
	Err   error
}

const prefixgt = "tags"

func GetTags(opts GetTagsOptions) (result *GetTagsResult) {
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	cacheKey := makeCacheKey(opts)
	if c, err := Cache.GetWithPrefix(prefixgt, cacheKey); err == nil {
		return c.(*GetTagsResult)
	}

	result = &GetTagsResult{Tags: []*modext.Tag{}}
	defer func() {
		if len(result.Tags) > 0 || result.Total > 0 || result.Err != nil {
			Cache.RemoveWithPrefix(prefixgt, cacheKey)
			Cache.SetWithPrefix(prefixgt, cacheKey, result, time.Hour*24*7)
		}
	}()

	q := []QueryMod{
		Select("tag.*", "COUNT(archive.tag_id) AS archive_count"),
		InnerJoin("archive_tags archive ON archive.tag_id = tag.id"),
		GroupBy("tag.id"),
		OrderBy("tag.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Tags(q...).BindG(context.Background(), &result.Tags)
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	count, err := models.Tags().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.ErrUnknown
		return
	}

	result.Total = int(count)
	return
}

func GetTagCount() (int64, error) {
	if c, err := Cache.Get("tag-count"); err == nil {
		return c.(int64), nil
	}

	count, err := models.Tags().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.ErrUnknown
	}

	Cache.Set("tag-count", count, time.Hour*24*7)
	return count, nil
}

var isTagValidMap = QueryMapCache{
	Map: make(map[string]bool),
}

func IsTagValid(str string) (isValid bool) {
	str = slugify(str)

	isTagValidMap.RLock()
	v, ok := isTagValidMap.Map[str]
	isTagValidMap.RUnlock()

	if ok {
		return v
	}

	result := GetTags(GetTagsOptions{})
	if result.Err != nil {
		return
	}

	defer func() {
		isTagValidMap.Lock()
		isTagValidMap.Map[str] = isValid
		isTagValidMap.Unlock()
	}()

	for _, tag := range result.Tags {
		if tag.Slug == str {
			isValid = true
			break
		}
	}
	return
}
