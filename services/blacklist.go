package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

type Blacklist struct {
	Archives map[string]bool
	Artists  map[string]bool

	once sync.Once
}

var blacklist Blacklist

func initBlacklist() {
	blacklist.once.Do(func() {
		blacklist.Archives = make(map[string]bool)
		blacklist.Artists = make(map[string]bool)

		if d, err := os.Stat(Config.Paths.Blacklist); os.IsNotExist(err) {
			return
		} else if d.IsDir() {
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

			t := strings.TrimSpace(arr[0])
			v := strings.TrimSpace(strings.Join(arr[1:], ":"))

			if t == "artist" {
				blacklist.Artists[v] = true
			} else if t == "title" {
				blacklist.Archives[v] = true
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
			return
		}
	})
}
