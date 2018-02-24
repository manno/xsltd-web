package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func handoff(w http.ResponseWriter, r *http.Request, requestPath string) {
	env := os.Environ()
	env = append(env, fmt.Sprintf("REQUEST_URI=%s", requestPath))
	env = append(env, fmt.Sprintf("DOCUMENT_ROOT=%s", config.WebRoot))
	env = append(env, fmt.Sprintf("HTTP_USER_AGENT=%s", r.UserAgent()))
	env = append(env, fmt.Sprintf("REMOTE_ADDR=%s", r.RemoteAddr))
	cmd := exec.Cmd{
		Args: []string{"xsltd-client.rb"},
		Env:  env,
		Path: config.ClientPath,
	}
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s", output)

}
