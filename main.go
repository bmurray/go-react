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
	proxy := flag.String("proxy", "http://localhost:3000/", "Address to proxy requests to")
	// Dir follows general filesystem pathing rules
	dir := flag.String("dir", "./react-app/build/", "Directory where the static built app resides")
	// Embed needs to be absolute, based on the arguments of pkger. See Makefile
	embed := flag.String("embed", "/react-app/build/", "Directory where the static built embedded app resides")
	mode := flag.String("mode", "proxy", "Mode to serve REACT site: proxy, dir, embed")
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
		// Dir mode is useful if you build your react app but don't want to embed it in the binary
		mux.Handle("/", http.FileServer(http.Dir(*dir)))
	case "embed":
		// Embed mode serves files that are embedded in the binary. Very useful for distribution
		mux.Handle("/", http.FileServer(EmbedDir(*embed)))
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
type EmbedDir string

func (d EmbedDir) Open(name string) (http.File, error) {
	// pkger.Dir can be used directly, but does not failback to the index file, breaking react-router
	pk := pkger.Dir(d)
	if f, err := pk.Open(name); err == nil {
		return f, err
	} else {
		return pk.Open("/index.html")
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
