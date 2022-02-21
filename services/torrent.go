package services

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "koushoku/config"

	"github.com/anacrolix/torrent"
)

var tClientCfg struct {
	*torrent.ClientConfig
	sync.Once
}

func DownloadBatches() {
	DownloadMagnets(Config.Paths.Batches)
}

func DownloadSingles() {
	DownloadMagnets(Config.Paths.Singles)
}

func newClient() (*torrent.Client, error) {
	tClientCfg.Once.Do(func() {
		tClientCfg.ClientConfig = torrent.NewDefaultClientConfig()
		tClientCfg.DataDir = Config.Directories.Data
		tClientCfg.DisableAggressiveUpload = true
		tClientCfg.HTTPUserAgent = "qBittorrent 4.4.1"
	})
	return torrent.NewClient(tClientCfg.ClientConfig)
}

func DownloadMagnets(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var magnets []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			magnets = append(magnets, line)
		}
	}
	f.Close()

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	client, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	total := len(magnets)
	log.Println(fmt.Sprintf("%d magnets found in %s", total, filepath.Base(file)))
	for i, magnet := range magnets {
		spec, err := torrent.TorrentSpecFromMagnetUri(magnet)
		if err != nil {
			log.Fatalln(err)
		}
		spec.DisableInitialPieceCheck = true

		t, _, err := client.AddTorrentSpec(spec)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(fmt.Sprintf("(%d/%d) Getting torrent info: %s", i+1, total, t.InfoHash().String()))
		<-t.GotInfo()

		info := t.Info()
		log.Println(
			fmt.Sprintf("(%d/%d) Downloading %s\n>> Hash: %s\n>> Files: %d | Size: %s",
				i+1,
				total,
				info.Name,
				t.InfoHash().String(),
				len(info.UpvertedFiles()),
				FormatBytes(info.TotalLength()),
			),
		)
		t.DownloadAll()
		t.Drop()
	}
	client.WaitAll()
}
