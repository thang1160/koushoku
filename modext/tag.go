package modext

import "koushoku/models"

type Tag struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count" boil:"archive_count"`
}

func NewTag(tag *models.Tag) *Tag {
	if tag == nil {
		return nil
	}
	return &Tag{
		ID:   tag.ID,
		Slug: tag.Slug,
		Name: tag.Name,
	}
}
