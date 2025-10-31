package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Config struct {
	Listen  string
	WebRoot string
	Xalan   string
}

func (config *Config) Setup() {
	config.Listen = fetch("LISTEN", "8080")
	config.WebRoot = fetch("WEBROOT", "/home/www/koeln.ccc.de/sandbox")
	config.Xalan = fetch("XALAN", "xalan")
}

func (config *Config) Validate() error {
	// Check if WebRoot exists
	if _, err := os.Stat(config.WebRoot); os.IsNotExist(err) {
		return fmt.Errorf("WEBROOT directory does not exist: %s", config.WebRoot)
	}

	// Check if xalan binary is available
	_, err := exec.LookPath(config.Xalan)
	if err != nil {
		return fmt.Errorf("xalan binary not found: %s (set XALAN environment variable to specify path)", config.Xalan)
	}

	return nil
}

func fetch(envKey string, defaultValue string) string {
	value := os.Getenv(envKey)
	if value != "" {
		return value
	}
	return defaultValue
}
