// Command line interface to hometube
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/janosgyerik/hometube"
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type params struct {
	url string
}

func parseArgs() params {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] url\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		exit()
	}

	url := flag.Args()[0]

	return params{
		url: url,
	}
}

func main() {
	downloader := hometube.DefaultDownloader()
	err := downloader.Init()
	if err != nil {
		log.Fatal(err)
	}

	p := parseArgs()

	err = downloader.Download(p.url, ".")
	if err != nil {
		log.Fatal(err)
	}
}
