package modext

import "koushoku/models"

type Artist struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count" boil:"archive_count"`
}

func NewArtist(artist *models.Artist) *Artist {
	if artist == nil {
		return nil
	}
	return &Artist{
		ID:   artist.ID,
		Slug: artist.Slug,
		Name: artist.Name,
	}
}
