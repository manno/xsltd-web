package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var config = new(Config)

func webrootPath(urlPath string) string {
	return path.Join(config.WebRoot, urlPath)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	absPath := webrootPath(r.URL.Path)
	requestPath := path.Clean(r.URL.Path)
	log.Println(r.RemoteAddr, requestPath, absPath)

	if fileExists(absPath) {
		if hasExtension(requestPath, "xml") {
			handoff(w, r, requestPath)
		} else {
			http.ServeFile(w, r, absPath)

		}
	} else {
		http.NotFound(w, r)
	}
}

func hasExtension(requestPath, extension string) bool {
	pos := strings.LastIndex(requestPath, extension)
	if pos > -1 && pos == len(requestPath)-len(extension) {
		if requestPath[pos-1] == '.' || requestPath[pos-1] == '?' {
			return true
		}
	}
	return false
}

func main() {
	config.Setup()

	var port = ":" + config.BindPort
	log.Printf("Listening on %s", config.BindPort)
	http.HandleFunc("/", viewHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
