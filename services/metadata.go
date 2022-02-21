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

	"koushoku/models"
	"koushoku/modext"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosimple/slug"

	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Metadata struct {
	Parody string
	Tags   []string
}

var wBaseURL string
var wUserAgent string
var wKeyword string

var wSelectorPrimay string
var wSelectorSecondary string
var wSelectorTertiary string
var wSelectorQuaternary string

func init() {
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

func ImportMetadata() {
	path := filepath.Join(Config.Directories.Root, "metadata.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	metadatas := make(map[string]*Metadata)
	if err := json.Unmarshal(buf, &metadatas); err != nil {
		log.Fatalln(err)
	}

	archives, err := models.Archives(
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Circle),
		Load(ArchiveRels.Magazine),
		Load(ArchiveRels.Parody),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	for _, arc := range archives {
		if arc.R != nil && (arc.R.Parody != nil || len(arc.R.Tags) > 0) {
			continue
		}

		fileName := filepath.Base(arc.Path)
		metadata, ok := metadatas[fileName]
		if !ok {
			continue
		}

		log.Println("Importing metadata of", fileName)
		archive := modext.NewArchive(arc).LoadRels(arc)
		archive.Parody = &modext.Parody{Name: metadata.Parody}
		archive.Tags = make([]*modext.Tag, len(metadata.Tags))
		for i, tag := range metadata.Tags {
			archive.Tags[i] = &modext.Tag{Name: tag}
		}

		if err := refreshArchiveRels(arc, archive); err != nil {
			log.Fatalln(err)
		}
	}
}

func ScrapeMetadata() {
	archives, err := models.Archives(
		Load(ArchiveRels.Parody),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	total := len(archives)
	log.Println(fmt.Sprintf("%d archives found", total))

	// Maximum number of concurrent requests
	c := make(chan bool, 10)
	defer func() {
		close(c)
	}()

	var wg sync.WaitGroup
	wg.Add(total)

	metadatas := make(map[string]*Metadata)
	for i, archive := range archives {
		if archive.R != nil && (archive.R.Parody != nil || len(archive.R.Tags) > 0) {
			continue
		}
		c <- true

		go func(i int, archive *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			if archive.R == nil || len(archive.R.Artists) == 0 {
				return
			}

			fileName := filepath.Base(archive.Path)
			log.Println(fileName) // DO NOT DELETE; intentional

			u := fmt.Sprintf("%s/search/%s", wBaseURL, archive.Slug)
			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				log.Fatalln(err)
			}

			req.Header.Set("User-Agent", wUserAgent)
			res, err := http.DefaultClient.Do(req)
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

			document, err := goquery.NewDocumentFromReader(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatalln(err)
			}

			var path string
			document.Find(wSelectorPrimay).Each(func(i int, s *goquery.Selection) {
				title := s.Find(wSelectorSecondary).First()
				if title == nil {
					return
				}

				if slug.Make(title.Text()) != archive.Slug {
					return
				}

				artist := s.Find(wSelectorTertiary).First()
				if artist == nil {
					return
				}

				slug := slug.Make(artist.Text())
				for _, artist := range archive.R.Artists {
					if slug == artist.Slug {
						path, _ = title.Attr("href")
						break
					}
				}
			})

			if len(path) == 0 {
				return
			}

			time.Sleep(time.Second)
			req, err = http.NewRequest("GET", wBaseURL+path, nil)
			if err != nil {
				log.Fatalln(err)
			}

			req.Header.Set("User-Agent", wUserAgent)
			res, err = http.DefaultClient.Do(req)
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

			metadata := &Metadata{}
			metadatas[fileName] = metadata

			fields := document.Find(wSelectorQuaternary)
			fields.EachWithBreak(func(i int, s *goquery.Selection) bool {
				if !strings.Contains(strings.ToLower(s.Text()), wKeyword) {
					return true
				}
				metadata.Parody = strings.TrimSpace(s.Children().Last().Text())
				return false
			})

			fields.Last().Children().First().Children().Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				if i > 0 && len(href) > 0 {
					metadata.Tags = append(metadata.Tags, strings.TrimSpace(s.Text()))
				}
			})
		}(i, archive)
	}

	wg.Wait()

	buf, err := json.Marshal(metadatas)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile("metadata.json", buf, 755); err != nil {
		log.Fatalln(err)
	}
}
