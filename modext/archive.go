package modext

import "koushoku/models"

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

	Artists   []*Artist   `json:"artists,omitempty"`
	Circles   []*Circle   `json:"circles,omitempty"`
	Magazines []*Magazine `json:"magazines,omitempty"`
	Parodies  []*Parody   `json:"parodies,omitempty"`
	Tags      []*Tag      `json:"tags,omitempty"`
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
	if model == nil || model.R == nil || len(model.R.Circles) == 0 {
		return archive
	}

	archive.Circles = make([]*Circle, len(model.R.Circles))
	for i, circle := range model.R.Circles {
		archive.Circles[i] = NewCircle(circle)
	}

	return archive
}

func (archive *Archive) LoadMagazine(model *models.Archive) *Archive {
	if model == nil || model.R == nil || len(model.R.Magazines) == 0 {
		return archive
	}

	archive.Magazines = make([]*Magazine, len(model.R.Magazines))
	for i, magazine := range model.R.Magazines {
		archive.Magazines[i] = NewMagazine(magazine)
	}

	return archive
}

func (archive *Archive) LoadParody(model *models.Archive) *Archive {
	if model == nil || model.R == nil || len(model.R.Parodies) == 0 {
		return archive
	}

	archive.Parodies = make([]*Parody, len(model.R.Parodies))
	for i, parody := range model.R.Parodies {
		archive.Parodies[i] = NewParody(parody)
	}

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
