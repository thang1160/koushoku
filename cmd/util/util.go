package main

import (
	"log"
	"os"

	"koushoku/database"
	"koushoku/services"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Delete    []int64 `long:"delete" description:"Delete archive(s) by id from the database"`
	Publish   []int64 `long:"publish" description:"Publish archive(s) by id"`
	Unpublish []int64 `long:"unpublish" description:"Unpublish archive(s) by id"`

	DeleteAll    bool `long:"delete-all" description:"Delete all archives from the database"`
	PublishAll   bool `long:"publish-all" description:"Publish all archives"`
	UnpublishAll bool `long:"unpublish-all" description:"Unpublish all archives"`

	Moderate        bool `long:"moderate" description:"Moderate all archives (blacklist)"`
	Purge           bool `long:"purge" description:"Purge symlinks"`
	PurgeThumbnails bool `long:"purge-thumbnails" description:"Purge thumbnails"`
	Index           bool `long:"index" description:"Index archives"`

	Scrape bool `long:"scrape" description:"Scrape metadata from you-know-where"`
	Import bool `long:"import" description:"Import metadata from metadata.json"`
}

func main() {
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		if !flags.WroteHelp(err) {
			log.Fatalln(err)
		}
		return
	}
	database.Init()

	if len(opts.Delete) > 0 {
		log.Println("Deleting archives from the database...")
		for _, id := range opts.Delete {
			if err := services.DeleteArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.DeleteAll {
		log.Println("Deleting all archives from the database...")
		if err := services.DeleteArchives(); err != nil {
			log.Fatalln(err)
		}
	}

	if opts.Purge {
		log.Println("Purging symlinks...")
		services.PurgeArchiveSymlinks()
	}

	if opts.PurgeThumbnails {
		log.Println("Purging thumbnails...")
		services.PurgeArchiveThumbnails()
	}

	if opts.Index {
		log.Println("Indexing archives...")
		services.IndexArchives()
	}

	if opts.Scrape {
		log.Println("Scraping metadata...")
		services.ScrapeMetadata()
	}

	if opts.Import {
		log.Println("Importing metadata...")
		services.ImportMetadata()
	}

	if opts.Moderate {
		log.Println("Moderating archives...")
		services.ModerateArchives()
	}

	if len(opts.Publish) > 0 {
		log.Println("Publishing archives...")
		for _, id := range opts.Publish {
			if _, err := services.PublishArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.PublishAll {
		log.Println("Publishing all archives...")
		if err := services.PublishArchives(); err != nil {
			log.Fatalln(err)
		}
	}

	if len(opts.Unpublish) > 0 {
		log.Println("Unpublishing archives...")
		for _, id := range opts.Unpublish {
			if _, err := services.UnpublishArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.UnpublishAll {
		log.Println("Unpublishing all archives...")
		if err := services.UnpublishArchives(); err != nil {
			log.Fatalln(err)
		}
	}
}
