package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	. "koushoku/config"

	"koushoku/database"
	"koushoku/models"
	"koushoku/modext"

	"github.com/PuerkitoBio/goquery"

	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type MetadataMap struct {
	Map  map[string]*Metadata
	once sync.Once
}

type Metadata struct {
	Parodies []string
	Tags     []string
}

var metadataMap MetadataMap

func initMetadata() {
	metadataMap.once.Do(func() {
		metadataMap.Map = make(map[string]*Metadata)
		path := filepath.Join(Config.Paths.Metadata)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return
		}

		buf, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.Unmarshal(buf, &metadataMap.Map); err != nil {
			log.Println(err)
		}
	})
}

func ImportMetadata() {
	initMetadata()

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
	defer func() {
		close(c)
	}()

	for _, arc := range archives {
		c <- true
		go func(arc *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			fileName := FileName(arc.Path)
			fileNameSlug := slugify(fileName)

			metadata, ok := metadataMap.Map[fileNameSlug]
			if !ok {
				return
			}

			log.Println("Importing metadata of", fileName)
			archive := modext.NewArchive(arc).LoadRels(arc)

			for _, parody := range metadata.Parodies {
				slug := slugify(parody)
				if v, ok := alias.Parodies[slug]; ok {
					slug = slugify(v)
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
				slug := slugify(tag)
				if v, ok := alias.Tags[slug]; ok {
					slug = slugify(v)
					tag = v
				}

				isDuplicate := false
				for _, tag2 := range archive.Tags {
					if tag2.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Tags = append(archive.Tags,
						&modext.Tag{Name: tag})
				}
			}

			if err := refreshArchiveRels(tx, arc, archive); err != nil {
				log.Fatalln(err)
			}
		}(arc)
	}
	wg.Wait()

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}
}

var wBaseURL string
var wUserAgent string
var wKeyword string

var wSelectorPrimay string
var wSelectorSecondary string
var wSelectorTertiary string
var wSelectorQuaternary string

func decodeStrings() {
	// Strings are encoded in base64 to avoid search engines
	doNotIndex := "WVVoU01HTklUVFpNZVRrelpETmpkVnB0Um5KaE0xVjFZbTFXTUE9PQ=="
	for !strings.HasPrefix(doNotIndex, "http") {
		buf, err := base64.StdEncoding.DecodeString(doNotIndex)
		if err != nil {
			log.Fatalln(err)
		}
		doNotIndex = string(buf)
	}
	wBaseURL = doNotIndex

	list := []string{
		"TW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVX",
		"ZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzk0LjAuNDYw",
		"Ni44MSBTYWZhcmkvNTM3LjM2",
	}
	buf, err := base64.StdEncoding.DecodeString(strings.Join(list, ""))
	if err != nil {
		log.Fatalln(err)
	}
	wUserAgent = string(buf)

	buf, err = base64.StdEncoding.DecodeString("cGFyb2R5")
	if err != nil {
		log.Fatalln(err)
	}
	wKeyword = string(buf)

	buf, err = base64.StdEncoding.DecodeString("Ym9keSA+IGRpdiAuZ3JpZCA+IGRpdltpZF49J2NvbnRlbnQtJ10=")
	if err != nil {
		log.Fatalln(err)
	}
	wSelectorPrimay = string(buf)

	buf, err = base64.StdEncoding.DecodeString("YS50ZXh0LWxn")
	if err != nil {
		log.Fatalln(err)
	}
	wSelectorSecondary = string(buf)

	buf, err = base64.StdEncoding.DecodeString("YS50ZXh0LXNt")
	if err != nil {
		log.Fatalln(err)
	}
	wSelectorTertiary = string(buf)

	list = []string{
		"Ym9keSA+IGRpdiA+IGRpdi5ncmlkID4gZGl2ID4gZGl2ID4gZGl2W2NsYXNzKj0n",
		"dGFibGUtY2VsbCddID4gLnRleHQtc20=",
	}
	buf, err = base64.StdEncoding.DecodeString(strings.Join(list, ""))
	if err != nil {
		log.Fatalln(err)
	}
	wSelectorQuaternary = string(buf)
}

func sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", wUserAgent)
	return http.DefaultClient.Do(req)
}

func ScrapeMetadata() {
	if len(wBaseURL) == 0 {
		decodeStrings()
	}

	initAlias()
	initMetadata()

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
	defer func() {
		close(c)
	}()

	var wg sync.WaitGroup
	wg.Add(total)

	for i, model := range archives {
		c <- true
		go func(i int, model *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			if model.R != nil && (len(model.R.Parodies) > 0 || len(model.R.Tags) > 0) {
				return
			}

			fileName := FileName(model.Path)
			fileNameSlug := slugify(fileName)

			if _, ok := metadataMap.Map[fileNameSlug]; ok {
				return
			}

			u := fmt.Sprintf("%s/search/%s", wBaseURL, model.Slug)
			res, err := sendRequest(u)
			if err != nil {
				log.Fatalln(err)
			}

			var path string
			var document *goquery.Document

			if res.StatusCode != http.StatusOK {
				if res.StatusCode == http.StatusNotFound {
					path = fmt.Sprintf("/hentai/%s-english", model.Slug)
					goto Skip
				}
				if res.StatusCode != http.StatusNotFound {
					log.Println("Failed to scrape metadata of", fileName)
				}
				res.Body.Close()
				return
			}

			document, err = goquery.NewDocumentFromReader(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatalln(err)
			}

			document.Find(wSelectorPrimay).Each(func(i int, s *goquery.Selection) {
				title := s.Find(wSelectorSecondary).First()
				if slugify(title.Text()) != model.Slug {
					return
				}

				str := slugify(s.Find(wSelectorTertiary).First().Text())
				if len(str) == 0 {
					return
				}

				if v, ok := alias.Artists[str]; ok {
					str = slugify(v)
				}

				for _, artist := range model.R.Artists {
					if str == artist.Slug {
						path, _ = title.Attr("href")
						break
					}
				}
			})

			if len(path) == 0 {
				return
			}

		Skip:
			time.Sleep(time.Second)
			res, err = sendRequest(wBaseURL + path)
			if err != nil {
				log.Fatalln(err)
			}

			if res.StatusCode != http.StatusOK {
				if res.StatusCode != http.StatusNotFound {
					log.Println("Failed to scrape metadata of", fileName)
				}
				res.Body.Close()
				return
			}

			document, err = goquery.NewDocumentFromReader(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("Found", fileName)

			metadata := &Metadata{}
			metadataMap.Map[fileNameSlug] = metadata

			fields := document.Find(wSelectorQuaternary)
			fields.EachWithBreak(func(i int, s *goquery.Selection) bool {
				if !strings.Contains(strings.ToLower(s.Text()), wKeyword) {
					return true
				}

				parodies := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
				for i := range parodies {
					parodies[i] = strings.TrimSpace(parodies[i])
					if v, ok := alias.Parodies[slugify(parodies[i])]; ok {
						parodies[i] = v
					}
				}

				metadata.Parodies = parodies
				return false
			})

			fields.Last().Children().First().Children().Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				if i > 0 && len(href) > 0 {
					tag := strings.TrimSpace(s.Text())
					if v, ok := alias.Tags[slugify(tag)]; ok {
						tag = v
					}
					metadata.Tags = append(metadata.Tags, tag)
				}
			})
		}(i, model)
	}
	wg.Wait()

	buf, err := json.Marshal(metadataMap.Map)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile("metadata.json", buf, 755); err != nil {
		log.Fatalln(err)
	}
}
