// Hometube server
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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
	items chan *item
}

func newQueue() *queue {
	return &queue{items: make(chan *item, 1)}
}

func worker(downloader hometube.Downloader, basedir string, q *queue) {
	for f := range q.items {
		log.Printf("worker: processing file: %s", f)
		err := downloader.Download(f.URL, basedir)
		if err != nil {
			log.Printf("failed to download %s; error: %s", f.URL, err)
		} else {
			log.Printf("successfully downloaded %s", f.URL)
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

type item struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type listDownloadedResponse struct {
	Items []item `json:"items"`
}

type server struct {
	queue   queue
	basedir string
}

// App holds template parameters for rendering the App's pages
type App struct {
	APIBaseURL string
}

func (s *server) download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	url := r.FormValue("url")
	f := &item{URL: url}
	fmt.Printf("adding to q: %s\n", url)
	s.queue.items <- f
	w.WriteHeader(http.StatusCreated)
	writeResponse(w, f)
}

func formatName(basedir string, fi os.FileInfo) string {
	if fi.Mode().IsDir() {
		files, err := ioutil.ReadDir(filepath.Join(basedir, fi.Name()))
		if err == nil {
			return fmt.Sprintf("%s (%d files)", fi.Name(), len(files))
		}
	}
	return fi.Name()
}

func (s *server) listDownloaded(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	osFiles, err := ioutil.ReadDir(s.basedir)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]item, len(osFiles))
	for index, f := range osFiles {
		items[index] = item{Name: formatName(s.basedir, f)}
	}
	response := listDownloadedResponse{Items: items}
	w.WriteHeader(http.StatusOK)
	writeResponse(w, response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("app/public/index.html")
	apiBaseURL := fmt.Sprintf("http://%s/api/v1", r.Host)
	t.Execute(w, &App{APIBaseURL: apiBaseURL})
}

func servePublicFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "app/public/"+r.RequestURI)
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
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.HandleFunc("/home", home).Methods(http.MethodGet)
	r.HandleFunc("/global.css", servePublicFile).Methods(http.MethodGet)
	r.HandleFunc("/build/bundle.css", servePublicFile).Methods(http.MethodGet)
	r.HandleFunc("/build/bundle.js", servePublicFile).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/download", s.download).
		Methods(http.MethodPost).
		Queries("url", "{url}")
	api.HandleFunc("/list/downloaded", s.listDownloaded).
		Methods(http.MethodGet)

	log.Printf("Listening on port %d, saving files to directory %s\n", args.port, args.basedir)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.port), r))
}
