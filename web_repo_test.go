package main

import (
	"testing"
)

func TestCleanPath(t *testing.T) {
	webRoot := "/home/www"
	var tests = []struct {
		input    string
		expected string
	}{
		{"meta", "/home/www/meta"},
		{"/some.xml", "/home/www/some.xml"},
		{"/path/some.xml", "/home/www/path/some.xml"},
		{"/path/../some.xml", "/home/www/some.xml"},
		{"/path/../../../some.xml", "/home/www/some.xml"},
		{"/path/some%20file.xml", "/home/www/path/some file.xml"},
	}

	repo := NewWebRepo(webRoot)
	for _, test := range tests {
		xmlFile, err := repo.cleanPath(test.input)
		if err != nil {
			t.Error("did not expect error")
		}
		got := xmlFile.FilesystemPath
		if got != test.expected {
			t.Errorf("expected %s to be %s, got: %s", test.input, test.expected, got)
		}
	}
}
