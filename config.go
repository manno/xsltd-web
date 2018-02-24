package main

import "os"

type Config struct {
	BindPort   string
	WebRoot    string
	ClientPath string
}

func (config *Config) Setup() {
	config.BindPort = fetch("PORT", "8080")
	config.WebRoot = fetch("WEBROOT", "/home/www/koeln.ccc.de")
	config.ClientPath = fetch("XSLTD_CLIENT", "/home/www/koeln.ccc.de/scripts/xsltd-client.rb")
}

func fetch(envKey string, defaultValue string) string {
	value := os.Getenv(envKey)
	if value != "" {
		return value
	}
	return defaultValue
}
