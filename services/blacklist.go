package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

var blacklist struct {
	Archives  map[string]bool
	ArchivesG []string
	Artists   map[string]bool
	Circles   map[string]bool
	Magazines map[string]bool
	Tags      map[string]bool

	once sync.Once
}

func initBlacklist() {
	blacklist.once.Do(func() {
		blacklist.Archives = make(map[string]bool)
		blacklist.Artists = make(map[string]bool)
		blacklist.Circles = make(map[string]bool)
		blacklist.Magazines = make(map[string]bool)
		blacklist.Tags = make(map[string]bool)

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

			arr := strings.Split(strings.ToLower(line), ":")
			if len(arr) < 2 {
				continue
			}

			v := slugify(strings.Join(arr[1:], ":"))
			switch strings.TrimSpace(arr[0]) {
			case "artist":
				blacklist.Artists[v] = true
			case "circle":
				blacklist.Circles[v] = true
			case "magazine":
				blacklist.Magazines[v] = true
			case "title":
				blacklist.Archives[v] = true
			case "title*":
				blacklist.ArchivesG = append(blacklist.ArchivesG, v)
			case "tag":
				blacklist.Tags[v] = true
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	})
}
