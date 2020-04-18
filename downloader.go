/*
Package hometube provides simple functions to download videos.
*/
package hometube

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path"
)

// Downloader is a common interface for utilities that download videos, after initialization.
type Downloader interface {
	Init() error
	Download(url, basedir string) error
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

type parsedURL struct {
	sanitized    string
	outputFormat string
}

func parseURL(s string) (*parsedURL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	m, _ := url.ParseQuery(u.RawQuery)

	listParam, hasListParam := m["list"]
	videoParam, hasVideoParam := m["v"]

	var param string
	var id string
	var outputFormat string

	if hasListParam {
		param = "list"
		id = listParam[0]
		outputFormat = "%(playlist)s/%(playlist_index)s - %(title)s.%(ext)s"
	} else if hasVideoParam {
		param = "v"
		id = videoParam[0]
		outputFormat = "%(title)s.%(ext)s"
	} else {
		return nil, errors.New("url must have v or list param")
	}

	return &parsedURL{
		fmt.Sprintf("https://www.youtube.com/watch?%s=%s", param, id),
		outputFormat,
	}, nil
}

// Download fetch the specified URL
func (youtubedl *YouTubeDl) Download(url, basedir string) error {
	pu, err := parseURL(url)
	if err != nil {
		return err
	}

	log.Printf("downloading %s ...\n", pu.sanitized)

	cmd := exec.Command(youtubedl.path, "-io", path.Join(basedir, pu.outputFormat), pu.sanitized)
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
