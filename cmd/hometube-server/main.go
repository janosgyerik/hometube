// Hometube server
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/janosgyerik/hometube"
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type params struct {
	basedir string
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func parseArgs() params {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] basedir\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		exit()
	}

	basedir := flag.Args()[0]
	if ok, _ := isDirectory(basedir); !ok {
		log.Fatalf("path does not exist or not a directory: %s", basedir)
	}

	return params{
		basedir: basedir,
	}
}

func writeResponse(w http.ResponseWriter, obj interface{}) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(bytes))
}

type message struct {
	Message string `json:"message"`
}

type file struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

type server struct {
	downloader hometube.Downloader
	basedir    string
}

func (s *server) download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := r.FormValue("url")
	filename := r.FormValue("filename")
	target := filepath.Join(s.basedir, filename)
	err := s.downloader.Download(url, target)
	if err != nil {
		writeResponse(w, &message{Message: err.Error()})
	} else {
		w.WriteHeader(http.StatusCreated)
		writeResponse(w, &file{URL: url, Filename: filename})
	}
}

func main() {
	downloader := hometube.DefaultDownloader()
	err := downloader.Init()
	if err != nil {
		log.Fatal(err)
	}

	args := parseArgs()
	s := &server{
		downloader: downloader,
		basedir:    args.basedir,
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/download", s.download).
		Methods(http.MethodPost).
		Queries("url", "{url}").
		Queries("filename", "{filename}")

	log.Fatal(http.ListenAndServe(":8080", r))
}
