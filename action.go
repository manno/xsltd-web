package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func handoff(w http.ResponseWriter, r *http.Request, xmlFile *XMLFile) {
	env := os.Environ()
	xslPath, err := NewChaosPage(config.WebRoot, xmlFile.FilesystemPath).Stylesheet()
	if err != nil {
		log.Printf("error extracting stylesheet: %s", err)
		http.Error(w, "Internal Server Error: failed to extract stylesheet", http.StatusInternalServerError)
		return
	}

	if xslPath == "" {
		log.Printf("no stylesheet found in: %s", xmlFile.FilesystemPath)
		http.Error(w, "Internal Server Error: no stylesheet declared in XML", http.StatusInternalServerError)
		return
	}

	cmd := exec.Cmd{
		Args: []string{"xalan", "-p", "thepath", "'" + xmlFile.URLPath + "'", xmlFile.FilesystemPath, xslPath},
		Env:  env,
		Path: config.Xalan,
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error running xalan: %s\n", err)
		log.Printf("command: %s %v", config.Xalan, cmd.Args)
		log.Printf("xalan output: %s", string(output))
		http.Error(w, "Internal Server Error: XSLT transformation failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", output)
}
