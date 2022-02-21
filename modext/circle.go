package modext

import "koushoku/models"

type Circle struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count" boil:"archive_count"`
}

func NewCircle(circle *models.Circle) *Circle {
	if circle == nil {
		return nil
	}
	return &Circle{
		ID:   circle.ID,
		Slug: circle.Slug,
		Name: circle.Name,
	}
}
