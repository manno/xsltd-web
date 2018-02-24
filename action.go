package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func handoff(w io.Writer, path string, requestPath string) {
	if strings.Index(path, config.WebRoot) != 0 {
		log.Fatal("directory listing outside webroot requested")
	}
	env := os.Environ()
	// env = append(env, fmt.Sprintf("MESSAGE_ID=%s", messageId))
	//   9 puts obj.handle(ENV['REQUEST_URI'], ENV['DOCUMENT_ROOT'], ENV['HTTP_USER_AGENT'], ENV['REMOTE_ADDR'])
	fmt.Printf("%#s\n", env)
	cmd := exec.Cmd{
		Args: []string{"xsltd-client.rb"},
		Env:  env,
		Path: config.ClientPath,
	}
	output, err := cmd.Output()
	if err != nil {
		w.Write(output)
	}
}
