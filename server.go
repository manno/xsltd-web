package main

import (
	"log"
	"net/http"
)

var config = new(Config)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	webPath := NewWebRepo(config.WebRoot)

	xmlFile, err := webPath.FindXML(r.URL.Path)
	if err != nil {
		log.Printf("failed to find xml for request: %s", err)
		return
	}

	log.Println(r.RemoteAddr, r.URL.Path, xmlFile.IsXML(), xmlFile.FilesystemPath)

	if !xmlFile.Exists() {
		http.NotFound(w, r)
	} else if !xmlFile.IsXML() {
		http.ServeFile(w, r, xmlFile.FilesystemPath)
	} else if xmlFile.IsXML() {
		handoff(w, r, xmlFile)
	}
}

func main() {
	config.Setup()

	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %s", err)
	}

	log.Printf("Listening on %s", config.Listen)
	log.Printf("Serving files from: %s", config.WebRoot)
	log.Printf("Using XSLT processor: %s", config.Xalan)

	http.HandleFunc("/", viewHandler)
	log.Fatal(http.ListenAndServe(config.Listen, nil))
}
