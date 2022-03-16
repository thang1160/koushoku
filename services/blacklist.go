package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

var blacklists struct {
	ArchiveMatches   map[string]bool
	ArchiveWildcards []string
	ArtistMatches    map[string]bool
	CircleMatches    map[string]bool
	MagazineMatches  map[string]bool
	TagMatches       map[string]bool

	once sync.Once
}

func loadBlacklists() {
	blacklists.once.Do(func() {
		blacklists.ArchiveMatches = make(map[string]bool)
		blacklists.ArtistMatches = make(map[string]bool)
		blacklists.CircleMatches = make(map[string]bool)
		blacklists.MagazineMatches = make(map[string]bool)
		blacklists.TagMatches = make(map[string]bool)

		stat, err := os.Stat(Config.Paths.Blacklist)
		if os.IsNotExist(err) || stat.IsDir() {
			return
		}

		f, err := os.Open(Config.Paths.Blacklist)
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
			if len(strs) < 2 {
				continue
			}

			v := Slugify(strings.Join(strs[1:], ":"))

			switch strings.TrimSpace(strs[0]) {
			case "title":
				blacklists.ArchiveMatches[v] = true
			case "title*":
				blacklists.ArchiveWildcards = append(blacklists.ArchiveWildcards, v)
			case "artist":
				blacklists.ArtistMatches[v] = true
			case "circle":
				blacklists.CircleMatches[v] = true
			case "magazine":
				blacklists.MagazineMatches[v] = true
			case "tag":
				blacklists.TagMatches[v] = true
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	})
}
