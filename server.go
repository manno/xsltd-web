package main

import (
	"log"
	"net/http"
	"os"
	"path"
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
		handoff(w, absPath, requestPath)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	config.Setup()

	var port = ":" + config.BindPort
	log.Printf("Listening on %s", config.BindPort)
	http.HandleFunc("/", viewHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
