package main

import (
	"log"
	"os"

	. "koushoku/config"

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

	Singles     bool   `long:"singles" description:"Download singles"`
	SinglesPath string `long:"singles-path" description:"Path to file containing single releases (optional)"`

	Batches     bool   `long:"batches" description:"Download batches"`
	BatchesPath string `long:"batches-path" description:"Path to file containing batch releases (optional)"`

	Moderate bool `long:"moderate" description:"Moderate all archives (blacklist)"`
	Purge    bool `long:"purge" description:"Purge symlinks"`
	Index    bool `long:"index" description:"Index archives"`

	ScrapeMetadata bool `long:"scrape-metadata" description:"Scrape metadata of all archives from you-know-where"`
	ImportMetadata bool `long:"import-metadata" description:"Import archives metadata from metadata.json"`
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

	if opts.Singles {
		log.Println("Downloading singles...")
		if len(opts.SinglesPath) > 0 {
			services.DownloadMagnets(opts.SinglesPath)
		} else {
			services.DownloadSingles()
		}
	}

	if opts.Batches {
		log.Println("Downloading batches...")
		if len(opts.BatchesPath) > 0 {
			services.DownloadMagnets(opts.BatchesPath)
		} else {
			services.DownloadBatches()
		}
	}

	if opts.Purge {
		log.Println("Purging symlinks...")
		services.PurgeArchiveSymlinks()
	}

	if opts.Index {
		log.Println("Indexing archives...")
		services.IndexArchives()
	}

	if opts.ScrapeMetadata {
		log.Println("Scraping metadata...")
		services.ScrapeMetadata()
	}

	if opts.ImportMetadata {
		log.Println("Importing metadata...")
		services.ImportMetadata()
	}

	if len(opts.Publish) > 0 {
		log.Println("Publishing archives...")
		for _, id := range opts.Publish {
			if _, err := services.PublishArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.Moderate {
		log.Println("Moderating archives...")
		services.ModerateArchives(Config.Paths.Blacklist)
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
