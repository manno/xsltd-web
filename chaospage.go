package main

import (
	"bufio"
	"os"
	"regexp"
)

type ChaosPage struct {
	XMLPath string
	WebRoot string
}

func NewChaosPage(webRoot string, xmlPath string) *ChaosPage {
	return &ChaosPage{WebRoot: webRoot, XMLPath: xmlPath}
}

func (c *ChaosPage) Stylesheet() (string, error) {
	r := regexp.MustCompile("<\\?xml-stylesheet")
	s := regexp.MustCompile("href=\"([^\"]+\\.xsl)\"")

	file, err := os.Open(c.XMLPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		if r.MatchString(scanner.Text()) {
			stylesheetName := s.FindStringSubmatch(scanner.Text())[1]
			p, err := prefixWebRoot(c.WebRoot, stylesheetName)
			if err != nil {
				return "", err
			}
			return p, nil
		}
	}
	return "", nil
}
