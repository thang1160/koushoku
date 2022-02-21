package modext

import "koushoku/models"

type Parody struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count" boil:"archive_count"`
}

func NewParody(parody *models.Parody) *Parody {
	if parody == nil {
		return nil
	}
	return &Parody{
		ID:   parody.ID,
		Slug: parody.Slug,
		Name: parody.Name,
	}
}
