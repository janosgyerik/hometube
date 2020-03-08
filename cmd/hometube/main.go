// Command line interface to hometube
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/janosgyerik/hometube"
)

const (
	// TODO placeholder for constants
	defaultCount = 5
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type params struct {
	url      string
	filename string
}

func parseArgs() params {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] url filename\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 2 {
		exit()
	}

	url := flag.Args()[0]
	filename := flag.Args()[1]

	return params{
		filename: filename,
		url:      url,
	}
}

func main() {
	downloader := hometube.DefaultDownloader()
	err := downloader.Init()
	if err != nil {
		log.Fatal(err)
	}

	p := parseArgs()

	err = downloader.Download(p.url, p.filename)
	if err != nil {
		log.Fatal(err)
	}
}
