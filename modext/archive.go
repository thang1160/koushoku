package modext

import (
	"koushoku/models"
)

type Archive struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`

	CreatedAt   int64 `json:"createdAt"`
	UpdatedAt   int64 `json:"updatedAt"`
	PublishedAt int64 `json:"publishedAt,omitempty"`

	Title string `json:"title"`
	Slug  string `json:"slug"`
	Pages int16  `json:"pages,omitempty"`
	Size  int64  `json:"size,omitempty"`

	Circle   *Circle   `json:"circle,omitempty"`
	Magazine *Magazine `json:"magazine,omitempty"`
	Parody   *Parody   `json:"parody,omitempty"`

	Artists []*Artist `json:"artists,omitempty"`
	Tags    []*Tag    `json:"tags,omitempty"`

	Thumbnail string `json:"thumbnail,omitempty"`
}

func NewArchive(model *models.Archive) *Archive {
	if model == nil {
		return nil
	}

	archive := &Archive{
		ID:   model.ID,
		Path: model.Path,

		CreatedAt: model.CreatedAt.Unix(),
		UpdatedAt: model.UpdatedAt.Unix(),

		Title: model.Title,
		Slug:  model.Slug,
		Pages: model.Pages,
		Size:  model.Size,
	}

	if model.PublishedAt.Valid {
		archive.PublishedAt = model.PublishedAt.Time.Unix()
	}

	return archive
}

func (archive *Archive) LoadRels(model *models.Archive) *Archive {
	if model == nil || model.R == nil {
		return archive
	}

	archive.LoadArtists(model)
	archive.LoadCircle(model)
	archive.LoadMagazine(model)
	archive.LoadParody(model)
	archive.LoadTags(model)

	return archive
}

func (archive *Archive) LoadArtists(model *models.Archive) *Archive {
	if model == nil || model.R == nil || len(model.R.Artists) == 0 {
		return archive
	}

	archive.Artists = make([]*Artist, len(model.R.Artists))
	for i, artist := range model.R.Artists {
		archive.Artists[i] = NewArtist(artist)
	}

	return archive
}

func (archive *Archive) LoadCircle(model *models.Archive) *Archive {
	if model == nil || model.R == nil || model.R.Circle == nil {
		return archive
	}

	archive.Circle = NewCircle(model.R.Circle)

	return archive
}

func (archive *Archive) LoadMagazine(model *models.Archive) *Archive {
	if model == nil || model.R == nil || model.R.Magazine == nil {
		return archive
	}

	archive.Magazine = NewMagazine(model.R.Magazine)

	return archive
}

func (archive *Archive) LoadParody(model *models.Archive) *Archive {
	if model == nil || model.R == nil || model.R.Parody == nil {
		return archive
	}

	archive.Parody = NewParody(model.R.Parody)

	return archive
}

func (archive *Archive) LoadTags(model *models.Archive) *Archive {
	if model == nil || model.R == nil || len(model.R.Tags) == 0 {
		return archive
	}

	archive.Tags = make([]*Tag, len(model.R.Tags))
	for i, tag := range model.R.Tags {
		archive.Tags[i] = NewTag(tag)
	}

	return archive
}
