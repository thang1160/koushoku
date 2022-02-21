package modext

import "koushoku/models"

type Magazine struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count" boil:"archive_count"`
}

func NewMagazine(magazine *models.Magazine) *Magazine {
	if magazine == nil {
		return nil
	}
	return &Magazine{
		ID:   magazine.ID,
		Slug: magazine.Slug,
		Name: magazine.Name,
	}
}
