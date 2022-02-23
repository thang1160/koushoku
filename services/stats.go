package services

import (
	"log"
	"math"
)

type Stats struct {
	ArchiveCount  int64
	ArtistCount   int64
	CircleCount   int64
	MagazineCount int64
	ParodyCount   int64
	TagCount      int64

	PageCount        int64
	AveragePageCount int64
	Size             int64
	AverageSize      int64
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

	if stats.ArchiveCount > 0 {
		if stats.PageCount > 0 {
			stats.AveragePageCount = int64(math.Round(float64(stats.PageCount) / float64(stats.ArchiveCount)))
		}
		if stats.Size > 0 {
			stats.AverageSize = int64(math.Round(float64(stats.Size) / float64(stats.ArchiveCount)))
		}
	}
	return
}

func GetStats() *Stats {
	return &stats
}
