package modext

import (
	"regexp"
	"strings"

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
	Size  string `json:"size,omitempty"`

	Circle   *Circle   `json:"circle,omitempty"`
	Magazine *Magazine `json:"magazine,omitempty"`
	Parody   *Parody   `json:"parody,omitempty"`

	Artists []*Artist `json:"artists,omitempty"`
	Tags    []*Tag    `json:"tags,omitempty"`

	Thumbnail string `json:"thumbnail,omitempty"`
}

func NewArchive(archive *models.Archive) *Archive {
	if archive == nil {
		return nil
	}

	c := &Archive{
		ID:   archive.ID,
		Path: archive.Path,

		CreatedAt: archive.CreatedAt.Unix(),
		UpdatedAt: archive.UpdatedAt.Unix(),

		Title: archive.Title,
		Slug:  archive.Slug,
		Pages: archive.Pages,
		Size:  archive.Size,
	}

	if archive.PublishedAt.Valid {
		c.PublishedAt = archive.PublishedAt.Time.Unix()
	}

	return c
}

func (c *Archive) LoadRels(archive *models.Archive) *Archive {
	if archive == nil || archive.R == nil {
		return c
	}

	c.LoadArtists(archive)
	c.LoadCircle(archive)
	c.LoadMagazine(archive)
	c.LoadParody(archive)
	c.LoadTags(archive)

	return c
}

func (c *Archive) LoadArtists(archive *models.Archive) *Archive {
	if archive == nil || archive.R == nil || len(archive.R.Artists) == 0 {
		return c
	}

	c.Artists = make([]*Artist, len(archive.R.Artists))
	for i, artist := range archive.R.Artists {
		c.Artists[i] = NewArtist(artist)
	}

	return c
}

func (c *Archive) LoadCircle(archive *models.Archive) *Archive {
	if archive == nil || archive.R == nil || archive.R.Circle == nil {
		return c
	}

	c.Circle = NewCircle(archive.R.Circle)

	return c
}

func (c *Archive) LoadMagazine(archive *models.Archive) *Archive {
	if archive == nil || archive.R == nil || archive.R.Magazine == nil {
		return c
	}

	c.Magazine = NewMagazine(archive.R.Magazine)

	return c
}

func (c *Archive) LoadParody(archive *models.Archive) *Archive {
	if archive == nil || archive.R == nil || archive.R.Parody == nil {
		return c
	}

	c.Parody = NewParody(archive.R.Parody)

	return c
}

func (c *Archive) LoadTags(archive *models.Archive) *Archive {
	if archive == nil || archive.R == nil || len(archive.R.Tags) == 0 {
		return c
	}

	c.Tags = make([]*Tag, len(archive.R.Tags))
	for i, tag := range archive.R.Tags {
		c.Tags[i] = NewTag(tag)
	}

	return c
}

var rgx = regexp.MustCompile(`(\(|\[)?[^\(\[\]\)]+(\)|\])?`)

func (c *Archive) FormatFromString(v string) {
	matches := rgx.FindAllString(v, -1)
	if len(matches) == 0 {
		return
	}

	var artists []string
	var circle, title, magazine string

	for i, match := range matches {
		match = strings.TrimSpace(match)
		if len(match) == 0 {
			continue
		}

		if i == 0 && strings.HasPrefix(match, "[") {
			match = strings.TrimPrefix(match, "[")
			match = strings.TrimSuffix(match, "]")

			artists = append(artists, match)
			continue
		}

		if i == 1 && strings.HasPrefix(match, "(") {
			match = strings.TrimPrefix(match, "(")
			match = strings.TrimSuffix(match, ")")

			if len(artists) > 0 {
				circle = artists[0]
				artists = artists[1:]
			}

			names := strings.Split(match, ",")
			for _, name := range names {
				artists = append(artists, strings.TrimSpace(name))
			}
			continue
		}

		if (i == 2 || i == 3) && strings.HasPrefix(match, "(") {
			match = strings.TrimPrefix(match, "(")
			match = strings.TrimSuffix(match, ")")

			if i < len(matches)-1 {
				next := matches[i+1]
				if len(next) > 0 && !(strings.HasPrefix(match, "(") || strings.HasSuffix(match, "[")) {
					continue
				}
			}

			if strings.HasPrefix(match, "x") || strings.EqualFold(match, "temp") ||
				strings.EqualFold(match, "strong") || strings.EqualFold(match, "complete") {
				continue
			}

			magazine = match
			continue
		}

		if i == 1 || i == 2 {
			title = strings.TrimSuffix(match, ".zip")
		}
	}

	c.Title = title
	if len(circle) > 0 {
		c.Circle = &Circle{Name: circle}
	}
	if len(magazine) > 0 {
		c.Magazine = &Magazine{Name: magazine}
	}

	c.Artists = make([]*Artist, len(artists))
	for i, artist := range artists {
		c.Artists[i] = &Artist{Name: artist}
	}
}
