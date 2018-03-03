package main

import "os"

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

func fetch(envKey string, defaultValue string) string {
	value := os.Getenv(envKey)
	if value != "" {
		return value
	}
	return defaultValue
}
