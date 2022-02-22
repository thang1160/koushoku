package services

import (
	"log"
)

type Stats struct {
	ArchiveCount  int64
	ArtistCount   int64
	CircleCount   int64
	MagazineCount int64
	ParodyCount   int64
	TagCount      int64

	PageCount        uint64
	AveragePageCount uint64
	Size             uint64
	AverageSize      uint64
}

var stats Stats

func AnalyzeStats() (err error) {
	log.Println("Analyzing stats...")
	defer func() {
		if err != nil {
			log.Println("AnalyzeStats returned an error:", err)
		}
	}()

	stats.ArchiveCount, err = GetArchiveCount()
	if err != nil {
		return
	}

	stats.Size, stats.PageCount, err = GetArchiveStats()
	if err != nil {
		return
	}

	stats.ArtistCount, err = GetArtistCount()
	if err != nil {
		return
	}

	stats.CircleCount, err = GetCircleCount()
	if err != nil {
		return
	}
	stats.MagazineCount, err = GetMagazineCount()
	if err != nil {
		return
	}

	stats.ParodyCount, err = GetParodyCount()
	if err != nil {
		return
	}

	stats.TagCount, err = GetTagCount()
	if err != nil {
		return
	}

	stats.AveragePageCount = stats.PageCount / uint64(stats.ArchiveCount)
	stats.AverageSize = stats.Size / uint64(stats.ArchiveCount)

	return
}

func GetStats() *Stats {
	return &stats
}
