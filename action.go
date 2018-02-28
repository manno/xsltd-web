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
		log.Println(err)
		return
	}

	cmd := exec.Cmd{
		//Args: []string{"xalan", "-q", "-param", "thepath", "'" + xmlFile.URLPath + "'", "-in", xmlFile.FilesystemPath, "-xsl", xslPath},
		Args: []string{"xalan", "-p", "thepath", "'" + xmlFile.URLPath + "'", xmlFile.FilesystemPath, xslPath},
		Env:  env,
		Path: config.Xalan,
	}
	output, err := cmd.Output()
	if err != nil {
		log.Printf("error running xalan: %s\n", err)
		log.Printf("%+v", []string{"xalan", "-p", "thepath", xmlFile.URLPath, xmlFile.FilesystemPath, xslPath})
	}
	fmt.Fprintf(w, "%s", output)

}
