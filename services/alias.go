package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

var alias struct {
	Archives  map[string]string
	Artists   map[string]string
	Circles   map[string]string
	Magazines map[string]string
	Parodies  map[string]string
	Tags      map[string]string

	once sync.Once
}

func initAlias() {
	alias.once.Do(func() {
		alias.Archives = make(map[string]string)
		alias.Artists = make(map[string]string)
		alias.Circles = make(map[string]string)
		alias.Magazines = make(map[string]string)
		alias.Parodies = make(map[string]string)
		alias.Tags = make(map[string]string)

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

			arr := strings.Split(strings.ToLower(line), ":")
			if len(arr) < 3 {
				continue
			}

			k := slugify(arr[1])
			v := strings.TrimSpace(strings.Join(arr[2:], ":"))

			switch strings.TrimSpace(arr[0]) {
			case "artist":
				alias.Artists[k] = v
			case "title":
				alias.Archives[k] = v
			case "circle":
				alias.Circles[k] = v
			case "magazine":
				alias.Magazines[k] = v
			case "parody":
				alias.Parodies[k] = v
			case "tag":
				alias.Tags[k] = v
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	})
}
