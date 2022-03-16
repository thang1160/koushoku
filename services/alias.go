package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

var aliases struct {
	ArchiveMatches  map[string]string
	ArtistMatches   map[string]string
	CircleMatches   map[string]string
	MagazineMatches map[string]string
	ParodyMatches   map[string]string
	TagMatches      map[string]string

	once sync.Once
}

func loadAliases() {
	aliases.once.Do(func() {
		aliases.ArchiveMatches = make(map[string]string)
		aliases.ArtistMatches = make(map[string]string)
		aliases.CircleMatches = make(map[string]string)
		aliases.MagazineMatches = make(map[string]string)
		aliases.ParodyMatches = make(map[string]string)
		aliases.TagMatches = make(map[string]string)

		stat, err := os.Stat(Config.Paths.Alias)
		if os.IsNotExist(err) || stat.IsDir() {
			return
		}

		f, err := os.Open(Config.Paths.Alias)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}

			strs := strings.Split(strings.ToLower(line), ":")
			if len(strs) < 3 {
				continue
			}

			k := Slugify(strs[1])
			v := strings.TrimSpace(strings.Join(strs[2:], ":"))

			switch strings.TrimSpace(strs[0]) {
			case "title":
				aliases.ArchiveMatches[k] = v
			case "artist":
				aliases.ArtistMatches[k] = v
			case "circle":
				aliases.CircleMatches[k] = v
			case "magazine":
				aliases.MagazineMatches[k] = v
			case "parody":
				aliases.ParodyMatches[k] = v
			case "tag":
				aliases.TagMatches[k] = v
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	})
}
