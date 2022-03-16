package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "koushoku/config"

	"koushoku/database"
	"koushoku/models"
	"koushoku/modext"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Metadata struct {
	Title     string
	Artists   []string
	Circles   []string
	Magazines []string
	Parodies  []string
	Tags      []string
}

var metadatas struct {
	Map  map[string]*Metadata
	once sync.Once
}

func loadMetadata() {
	metadatas.once.Do(func() {
		metadatas.Map = make(map[string]*Metadata)
		path := filepath.Join(Config.Paths.Metadata)

		stat, err := os.Stat(path)
		if os.IsNotExist(err) || stat.IsDir() {
			return
		}

		buf, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.Unmarshal(buf, &metadatas.Map); err != nil {
			log.Println(err)
		}
	})
}

func ImportMetadata() {
	loadMetadata()

	archives, err := models.Archives(
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Circles),
		Load(ArchiveRels.Magazines),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := database.Conn.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(archives))

	c := make(chan bool, 20)
	defer close(c)

	for _, model := range archives {
		c <- true
		go func(model *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			fn := FileName(model.Path)
			fnSlug := Slugify(fn)

			metadata, ok := metadatas.Map[fnSlug]
			if !ok {
				return
			}

			log.Println("Importing metadata of", fn)
			archive := modext.NewArchive(model).LoadRels(model)

			for _, artist := range metadata.Artists {
				slug := Slugify(artist)
				if v, ok := aliases.ArtistMatches[slug]; ok {
					slug = Slugify(v)
					artist = v
				}

				isDuplicate := false
				for _, a := range archive.Artists {
					if a.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Artists = append(archive.Artists,
						&modext.Artist{Name: artist})
				}
			}

			for _, circle := range metadata.Circles {
				slug := Slugify(circle)
				if v, ok := aliases.CircleMatches[slug]; ok {
					slug = Slugify(v)
					circle = v
				}

				isDuplicate := false
				for _, c := range archive.Circles {
					if c.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Circles = append(archive.Circles,
						&modext.Circle{Name: circle})
				}
			}

			for _, magazine := range metadata.Magazines {
				slug := Slugify(magazine)
				if v, ok := aliases.MagazineMatches[slug]; ok {
					slug = Slugify(v)
					magazine = v
				}

				isDuplicate := false
				for _, m := range archive.Magazines {
					if m.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Magazines = append(archive.Magazines,
						&modext.Magazine{Name: magazine})
				}
			}

			for _, parody := range metadata.Parodies {
				slug := Slugify(parody)
				if v, ok := aliases.ParodyMatches[slug]; ok {
					slug = Slugify(v)
					parody = v
				}

				isDuplicate := false
				for _, p := range archive.Parodies {
					if p.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Parodies = append(archive.Parodies,
						&modext.Parody{Name: parody})
				}
			}

			for _, tag := range metadata.Tags {
				slug := Slugify(tag)
				if v, ok := aliases.TagMatches[slug]; ok {
					slug = Slugify(v)
					tag = v
				}

				isDuplicate := false
				for _, t := range archive.Tags {
					if t.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Tags = append(archive.Tags,
						&modext.Tag{Name: tag})
				}
			}

			if len(metadata.Title) > 0 && metadata.Title != archive.Title {
				model.Title = metadata.Title
				model.Slug = Slugify(model.Title)

				if v, ok := aliases.ArchiveMatches[model.Slug]; ok {
					model.Slug = Slugify(v)
					model.Title = v
				}

				if err := model.Update(tx, boil.Whitelist(ArchiveCols.Title, ArchiveCols.Slug)); err != nil {
					log.Fatalln(err)
				}
			}

			if err := populateArchiveRels(tx, model, archive); err != nil {
				log.Fatalln(err)
			}
		}(model)
	}
	wg.Wait()

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}
}

const fBaseURL = "https://www.fakku.net"
const iBaseURL = "https://irodoricomics.com"

var httpClient struct {
	*http.Client
	once sync.Once
}

func initHttpClient() {
	httpClient.once.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalln(err)
		}

		u, err := url.Parse("https://fakku.net")
		if err != nil {
			log.Fatalln(err)
		}

		httpClient.Client = &http.Client{Jar: jar}
		jar.SetCookies(u, []*http.Cookie{{
			Name:     "fakku_sid",
			Value:    Config.HTTP.Cookie,
			Domain:   "fakku.net",
			HttpOnly: true,
		}})
	})
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36"

func sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	return httpClient.Do(req)
}

func searchF(model *models.Archive) (path string, err error) {
	var (
		res      *http.Response
		document *goquery.Document
	)

	res, err = sendRequest(fmt.Sprintf("%s/search/%s", fBaseURL, model.Slug))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	document, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	document.Find("body > div .grid > div[id^='content-']").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			titleElement := s.Find("a.text-lg").First()
			if Slugify(titleElement.Text()) != model.Slug {
				return true
			}

			artistSlug := Slugify(s.Find("a.text-sm").First().Text())
			if len(artistSlug) == 0 {
				return true
			}

			if v, ok := aliases.ArtistMatches[artistSlug]; ok {
				artistSlug = Slugify(v)
			}

			for _, artist := range model.R.Artists {
				if artistSlug == artist.Slug {
					path, _ = titleElement.Attr("href")
					break
				}
			}

			return len(path) == 0
		})
	return
}

func scrapeF(fn, fnSlug string, model *models.Archive) (ok bool) {
	path, err := searchF(model)
	if err != nil {
		log.Fatalln(err)
	}

	if len(path) == 0 {
		path = fmt.Sprintf("/hentai/%s-english", model.Slug)
	}

	res, err := sendRequest(fBaseURL + path)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[F] metadata not available:", fn)
		return
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("[F] metadata found:", fn)

	metadata, ok := metadatas.Map[fnSlug]
	if !ok {
		metadata = &Metadata{}
		metadatas.Map[fnSlug] = metadata
	}

	metadata.Title = strings.TrimSpace(document.Find("body > div > div.grid > div > div > div[class*='table-cell'] > h1").Text())

	fields := document.Find("body > div > div.grid > div > div > div[class*='table-cell'] > .text-sm")
	fields.Each(func(i int, s *goquery.Selection) {
		if s.Children().Length() == 1 {
			return
		}

		section := strings.ToLower(s.Children().First().Text())
		if strings.Contains(section, "artist") {
			artists := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, artist := range artists {
				artist = strings.TrimSpace(artist)
				if v, ok := aliases.ArtistMatches[Slugify(artist)]; ok {
					artist = v
				}

				duplicate := false
				for _, v := range metadata.Artists {
					if v == artist {
						duplicate = true
						break
					}
				}

				if !duplicate {
					metadata.Artists = append(metadata.Artists, artist)
				}
			}
		} else if strings.Contains(section, "circle") {
			circles := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, circle := range circles {
				circle = strings.TrimSpace(circle)
				if v, ok := aliases.CircleMatches[Slugify(circle)]; ok {
					circle = v
				}

				duplicate := false
				for _, v := range metadata.Circles {
					if v == circle {
						duplicate = true
						break
					}
				}
				if !duplicate {
					metadata.Circles = append(metadata.Circles, circle)
				}
			}
		} else if strings.Contains(section, "parody") {
			parodies := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, parody := range parodies {
				parody = strings.TrimSpace(parody)
				if v, ok := aliases.ParodyMatches[Slugify(parody)]; ok {
					parody = v
				}

				duplicate := false
				for _, v := range metadata.Parodies {
					if v == parody {
						duplicate = true
						break
					}
				}

				if !duplicate {
					metadata.Parodies = append(metadata.Parodies, parody)
				}
			}
		} else if strings.Contains(section, "magazine") {
			magazines := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, magazine := range magazines {
				magazine = strings.TrimSpace(magazine)
				if v, ok := aliases.MagazineMatches[Slugify(magazine)]; ok {
					magazine = v
				}

				duplicate := false
				for _, v := range metadata.Magazines {
					if v == magazine {
						duplicate = true
						break
					}
				}

				if !duplicate {
					metadata.Magazines = append(metadata.Magazines, magazine)
				}
			}
		}
	})

	// Parse tags
	fields.Last().Children().First().Children().Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if i > 0 && len(href) > 0 {
			tag := strings.TrimSpace(s.Text())
			if v, ok := aliases.TagMatches[Slugify(tag)]; ok {
				tag = v
			}

			duplicate := false
			for _, v := range metadata.Tags {
				if v == tag {
					duplicate = true
					break
				}
			}

			if !duplicate {
				metadata.Tags = append(metadata.Tags, tag)
			}
		}
	})
	return true
}

func searchI(model *models.Archive) (path string, err error) {
	var (
		res      *http.Response
		document *goquery.Document
	)

	res, err = sendRequest(fmt.Sprintf("%s/index.php?route=product/search&search=%s", iBaseURL, model.Title))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	document, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
	}

	entries := document.Find(".main-products > .product-layout")
	entries.EachWithBreak(func(i int, s *goquery.Selection) bool {
		titleElement := s.Find(".caption > .name a")
		if Slugify(titleElement.Text()) == model.Slug {
			return true
		}

		artistSlug := Slugify(s.Find(".caption > .stats span a").Text())
		if len(artistSlug) == 0 {
			return true
		}

		for _, artist := range model.R.Artists {
			if artistSlug == artist.Slug {
				path, _ = titleElement.Attr("href")
				break
			}
		}

		return len(path) == 0
	})

	if len(path) == 0 {
		document.EachWithBreak(func(i int, s *goquery.Selection) bool {
			titleElement := s.Find(".caption > .name a")
			if strings.Contains(Slugify(titleElement.Text()), model.Slug) {
				return true
			}

			artistSlug := Slugify(s.Find(".caption > .stats span a").Text())
			if len(artistSlug) == 0 {
				return true
			}

			for _, artist := range model.R.Artists {
				if artistSlug == artist.Slug {
					path, _ = titleElement.Attr("href")
					break
				}
			}

			return len(path) == 0
		})
	}
	return
}

func scrapeI(fn, fnSlug string, model *models.Archive) (ok bool) {
	path, err := searchI(model)
	if err != nil {
		log.Fatalln(err)
	}

	if len(path) == 0 && len(model.R.Artists) == 1 {
		path = fmt.Sprintf("/%s/%s", model.R.Artists[0].Slug, model.Slug)
	}

	res, err := sendRequest(iBaseURL + path)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[I] metadata not available:", fn)
		return
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Print("[I] metadata found:", fn)

	metadata, ok := metadatas.Map[fnSlug]
	if !ok {
		metadata = &Metadata{}
		metadatas.Map[fnSlug] = metadata
	}

	metadata.Title = strings.TrimSpace(document.Find("h1.title.page-title").Text())

	artists := document.Find(".product-manufacturer a")
	if artists.Length() > 0 {
		artists.Each(func(i int, s *goquery.Selection) {
			artist := strings.TrimSpace(s.Text())
			if v, ok := aliases.ArtistMatches[Slugify(artist)]; ok {
				artist = v
			}

			duplicate := false
			for _, v := range metadata.Artists {
				if v == artist {
					duplicate = true
					break
				}
			}

			if !duplicate {
				metadata.Artists = append(metadata.Artists, artist)
			}
		})
	}

	tags := document.Find(".ctags a")
	if tags.Length() > 0 {
		tags.Each(func(i int, s *goquery.Selection) {
			tag := strings.TrimSpace(s.Text())
			if v, ok := aliases.TagMatches[Slugify(tag)]; ok {
				tag = v
			}

			duplicate := false
			for _, v := range metadata.Tags {
				if v == tag {
					duplicate = true
					break
				}
			}

			if !duplicate {
				metadata.Tags = append(metadata.Tags, tag)
			}
		})
	}
	return true
}

func ScrapeMetadata() {
	initHttpClient()
	loadAliases()
	loadMetadata()

	archives, err := models.Archives(
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	total := len(archives)
	log.Println(fmt.Sprintf("%d archives found", total))

	c := make(chan bool, 10)
	defer close(c)

	var wg sync.WaitGroup
	wg.Add(total)

	for i, model := range archives {
		c <- true
		go func(i int, model *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			fn := FileName(model.Path)
			fnSlug := Slugify(fn)

			if _, ok := metadatas.Map[fnSlug]; ok {
				return
			}

			if !scrapeF(fn, fnSlug, model) {
				scrapeI(fn, fnSlug, model)
			}
		}(i, model)
	}
	wg.Wait()

	buf, err := json.Marshal(metadatas.Map)
	if err == nil {
		err = os.WriteFile("metadata.json", buf, 755)
	}

	if err != nil {
		log.Fatalln(errors.WithStack(err))
	}
}
