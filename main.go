package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/markbates/pkger"
)

func main() {
	mode := flag.String("mode", "proxy", "Mode to serve REACT site: proxy, dir, embed")
	// Proxy mode proxies all other connections to the npm server
	proxy := flag.String("proxy", "http://localhost:3000/", "Address to proxy requests to")
	// Dir follows general filesystem pathing rules
	dir := flag.String("dir", "./react-app/build/", "Directory where the static built app resides")
	// Embed needs to be absolute, based on the arguments of pkger. See Makefile
	embed := flag.String("embed", "/react-app/build/", "Directory where the static built embeded app resides")
	listen := flag.String("listen", ":8080", "Listen on address")
	flag.Parse()

	// Basic ServeMux and API that just sends the time
	mux := http.NewServeMux()
	mux.HandleFunc("/api", basicAPI)

	// The React serve magic
	switch *mode {
	case "proxy":
		// Proxy mode is most useful for development
		// Preserves live-reload
		u, err := url.Parse(*proxy)
		if err != nil {
			log.Fatalf("Cannot parse proxy address: %s", err)
		}
		mux.Handle("/", httputil.NewSingleHostReverseProxy(u))
	case "dir":
		// Dir mode is useful if you build your react app but don't want to embed it in the binary, such as Docker deploys
		mux.Handle("/", http.FileServer(EmbedDir{http.Dir(*dir)}))
	case "embed":
		// Embed mode serves files that are embedded in the binary. Very useful for one-file distribution
		mux.Handle("/", http.FileServer(EmbedDir{pkger.Dir(*embed)}))
	default:
		// Any other mode would assume you have a reverse proxy, like nginx, that filters traffic
		log.Println("No react mode; this only works if you have a frontend reverse proxy")
	}
	s := &http.Server{
		Addr:         *listen,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println(s.ListenAndServe())
}

// EmbedDir provides a convenience method to default requests back to /index.html, allowing react-router to work correctly
type EmbedDir struct {
	http.FileSystem
}

// Open implementation of http.FileSystem that falls back to serving /index.html, allowing react-router to operate
func (d EmbedDir) Open(name string) (http.File, error) {
	if f, err := d.FileSystem.Open(name); err == nil {
		return f, err
	} else {
		return d.FileSystem.Open("/index.html")
	}
}

func basicAPI(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	j := json.NewEncoder(w)
	if err := j.Encode(struct {
		Time string `json:"time"`
	}{
		time.Now().String(),
	}); err != nil {
		log.Println("Cannot encode time", err)
	}
}
