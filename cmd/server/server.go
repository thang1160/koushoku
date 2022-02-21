package main

import (
	"log"
	"os"

	. "koushoku/config"

	"koushoku/controllers"
	"koushoku/database"
	"koushoku/server"
)

func init() {
	os.Setenv("MALLOC_ARENA_MAX", "2")
}

func main() {
	if _, err := os.Stat(Config.Directories.Data); os.IsNotExist(err) {
		if err := os.MkdirAll(Config.Directories.Data, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := os.Stat(Config.Directories.Thumbnails); os.IsNotExist(err) {
		if err := os.MkdirAll(Config.Directories.Thumbnails, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	database.Init()
	server.Init()
	controllers.Init()
	server.Start()
}
