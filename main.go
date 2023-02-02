package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	down "github.com/SantiagoBedoya/down/pkg"
)

func main() {
	var (
		url     string
		dest    string
		mode    int
		workers int
	)
	home, _ := os.UserHomeDir()
	flag.StringVar(&url, "url", "", "File URL to Download")
	flag.StringVar(&dest, "dest", home, "Destination folder")
	flag.IntVar(&mode, "mode", 0, "Download Mode (concurrent: 0 | normal: 1)")
	flag.IntVar(&workers, "workers", 5, "Worker for concurrent download")

	flag.Parse()

	if url == "" {
		log.Println("url should not be empty")
		flag.PrintDefaults()
		os.Exit(1)
	}

	destPath := getDestinationPath(url, dest)

	app := &down.App{
		Concurrency: workers,
		URI:         url,
		Chunks:      make(map[int][]byte),
		Err:         nil,
		Destination: destPath,
		Mutex:       &sync.Mutex{},
	}

	switch mode {
	case 0:
		if err := app.Concurrent(); err != nil {
			log.Printf("Error downloading file in concurrent mode: %v\n", err)
			os.Exit(1)
		}
	case 1:
		if err := app.Normal(); err != nil {
			log.Printf("Error downloading file in normal mode: %v\n", err)
			os.Exit(1)
		}
	default:
		log.Println("Unknown mode")
		os.Exit(1)
	}

}

func getDestinationPath(url, dest string) string {
	urlParts := strings.Split(url, "/")
	fileName := urlParts[len(urlParts)-1]
	destPath := fmt.Sprintf("%s/%s", dest, fileName)
	return destPath
}
