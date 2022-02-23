package main

import (
	"os"

	"koushoku/controllers"
	"koushoku/database"
	"koushoku/server"
	"koushoku/services"
)

func init() {
	os.Setenv("MALLOC_ARENA_MAX", "2")
}

func main() {
	database.Init()
	if err := services.AnalyzeStats(); err != nil {
		return
	}

	server.Init()
	controllers.Init()
	server.Start()
}
