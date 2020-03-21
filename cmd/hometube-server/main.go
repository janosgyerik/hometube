// Hometube server
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/janosgyerik/hometube"
)

const (
	defaultBasedir = "."
	defaultPort    = 8080
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type params struct {
	basedir string
	port    int
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
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	portPtr := flag.Int("port", defaultPort, "the port to listen on")
	basedir := flag.String("basedir", defaultBasedir, "the base directory to download files to")

	flag.Parse()

	if flag.NArg() > 0 {
		exit()
	}

	if ok, _ := isDirectory(*basedir); !ok {
		log.Fatalf("path does not exist or not a directory: %s", *basedir)
	}

	return params{
		basedir: *basedir,
		port:    *portPtr,
	}
}

type queue struct {
	items chan *file
}

func newQueue() *queue {
	return &queue{items: make(chan *file, 1)}
}

func worker(downloader hometube.Downloader, basedir string, q *queue) {
	for f := range q.items {
		log.Printf("worker: processing file: %s", f)
		target := filepath.Join(basedir, f.Filename)
		err := downloader.Download(f.URL, target)
		if err != nil {
			log.Printf("failed to download %s to %s", f.URL, target)
		} else {
			log.Printf("successfully downloaded %s to %s", f.URL, target)
		}
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

type listDownloadedResponse struct {
	Files []file `json:"files"`
}

type server struct {
	queue   queue
	basedir string
}

func (s *server) download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	url := r.FormValue("url")
	filename := r.FormValue("filename")
	f := &file{URL: url, Filename: filename}
	s.queue.items <- f
	w.WriteHeader(http.StatusCreated)
	writeResponse(w, f)
}

func (s *server) listDownloaded(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	osFiles, err := ioutil.ReadDir(s.basedir)
	if err != nil {
		log.Fatal(err)
	}

	files := make([]file, 0)
	for _, f := range osFiles {
		files = append(files, file{Filename: f.Name()})
	}
	response := listDownloadedResponse{Files: files}
	w.WriteHeader(http.StatusOK)
	writeResponse(w, response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	// TODO bundle the svelte app
	// TODO serve the single file / template from svelte app
}

func main() {
	downloader := hometube.DefaultDownloader()
	err := downloader.Init()
	if err != nil {
		log.Fatal(err)
	}

	args := parseArgs()

	q := newQueue()

	go worker(downloader, args.basedir, q)

	s := &server{queue: *q, basedir: args.basedir}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/home", home).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/download", s.download).
		Methods(http.MethodPost).
		Queries("url", "{url}").
		Queries("filename", "{filename}")
	api.HandleFunc("/list/downloaded", s.listDownloaded).
		Methods(http.MethodGet)

	log.Printf("Listening on port %d, saving files to directory %s\n", args.port, args.basedir)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.port), r))
}
