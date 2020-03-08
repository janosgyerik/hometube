/*
Package hometube provides simple functions to download videos.
*/
package hometube

import (
	"bufio"
	"log"
	"os/exec"
)

// Downloader is a common interface for utilities that download videos, after initialization.
type Downloader interface {
	Init() error
	Download(url, filename string) error
}

// YouTubeDl implements Downloader using youtube-dl
type YouTubeDl struct {
	path string
}

// Init checks if the youtube-dl command exists on PATH
func (youtubedl *YouTubeDl) Init() error {
	path, err := exec.LookPath("youtube-dl")
	if err == nil {
		youtubedl.path = path
	}
	return err
}

// Download downloads the specified URL and saves to the specified filename
func (youtubedl *YouTubeDl) Download(url, filename string) error {
	log.Printf("downloading %s as %s ...\n", url, filename)

	cmd := exec.Command(youtubedl.path, "-o", filename, url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				return
			}
			log.Printf("youtube-dl: %s", line)
		}
	}()
	go func() {
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				return
			}
			log.Printf("ERROR: %s", line)
		}
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}
	return err
}

// DefaultDownloader creates and returns a downloader utility
func DefaultDownloader() Downloader {
	return &YouTubeDl{}
}
