package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestStylesheet(t *testing.T) {
	webRoot := "/home/www"

	var tests = []struct {
		content  string
		expected string
	}{
		{`<?xml-stylesheet href="/xsl/page_sub.xsl" type="text/xsl"?>`, "/xsl/page_sub.xsl"},
		{`<?xml-stylesheet href="xsl/page_sub.xsl" type="text/xsl"?>`, "/xsl/page_sub.xsl"},
		{`<?xml-stylesheet href="/xsl/page_teaser.xsl" type="text/xsl"?>`, "/xsl/page_teaser.xsl"},
		{`<?xml-stylesheet href="xsl/page_index.xsl" type="text/xsl"?>`, "/xsl/page_index.xsl"},
	}

	for _, test := range tests {
		tmpfile := NewTempFile(test.content)
		defer os.Remove(tmpfile.Name())

		c := NewChaosPage(webRoot, tmpfile.Name())

		s, err := c.Stylesheet()
		if err != nil {
			t.Error("did not expect error")
		}
		if s != webRoot+test.expected {
			t.Errorf("expected %s, got: %s", webRoot+test.expected, s)
		}
		os.Remove(tmpfile.Name())
	}
}

func NewTempFile(content string) *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.WriteString(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile
}
