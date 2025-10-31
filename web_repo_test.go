package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestFindXML(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "webroot-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	indexXML := filepath.Join(tmpDir, "index.xml")
	pageXML := filepath.Join(tmpDir, "page.xml")
	ioutil.WriteFile(indexXML, []byte("<root/>"), 0644)
	ioutil.WriteFile(pageXML, []byte("<root/>"), 0644)

	repo := NewWebRepo(tmpDir)

	t.Run("empty path treated as current directory", func(t *testing.T) {
		xmlFile, err := repo.FindXML("")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		// Empty string gets cleaned to "." by path.Clean
		if xmlFile.URLPath != "." {
			t.Errorf("expected URLPath '.', got '%s'", xmlFile.URLPath)
		}
	})

	t.Run("root path returns index.xml", func(t *testing.T) {
		xmlFile, err := repo.FindXML("/")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if xmlFile.URLPath != "index.xml" {
			t.Errorf("expected URLPath 'index.xml', got '%s'", xmlFile.URLPath)
		}
		if xmlFile.FilesystemPath != indexXML {
			t.Errorf("expected FilesystemPath '%s', got '%s'", indexXML, xmlFile.FilesystemPath)
		}
	})

	t.Run("exact xml file match", func(t *testing.T) {
		xmlFile, err := repo.FindXML("/page.xml")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if xmlFile.URLPath != "/page.xml" {
			t.Errorf("expected URLPath '/page.xml', got '%s'", xmlFile.URLPath)
		}
		if xmlFile.FilesystemPath != pageXML {
			t.Errorf("expected FilesystemPath '%s', got '%s'", pageXML, xmlFile.FilesystemPath)
		}
	})

	t.Run("html request maps to xml file", func(t *testing.T) {
		xmlFile, err := repo.FindXML("/page.html")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if xmlFile.URLPath != "/page.xml" {
			t.Errorf("expected URLPath '/page.xml', got '%s'", xmlFile.URLPath)
		}
		if xmlFile.FilesystemPath != pageXML {
			t.Errorf("expected FilesystemPath '%s', got '%s'", pageXML, xmlFile.FilesystemPath)
		}
	})

	t.Run("path without extension appends .xml", func(t *testing.T) {
		xmlFile, err := repo.FindXML("/page")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if xmlFile.URLPath != "/page.xml" {
			t.Errorf("expected URLPath '/page.xml', got '%s'", xmlFile.URLPath)
		}
		if xmlFile.FilesystemPath != pageXML {
			t.Errorf("expected FilesystemPath '%s', got '%s'", pageXML, xmlFile.FilesystemPath)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		xmlFile, err := repo.FindXML("/nonexistent.xml")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		expectedPath := filepath.Join(tmpDir, "nonexistent.xml")
		if xmlFile.FilesystemPath != expectedPath {
			t.Errorf("expected FilesystemPath '%s', got '%s'", expectedPath, xmlFile.FilesystemPath)
		}
	})
}

func TestXMLFileHelpers(t *testing.T) {
	t.Run("IsXML identifies xml files", func(t *testing.T) {
		xmlFile := &XMLFile{FilesystemPath: "/path/to/file.xml"}
		if !xmlFile.IsXML() {
			t.Error("expected IsXML() to return true for .xml file")
		}

		nonXMLFile := &XMLFile{FilesystemPath: "/path/to/file.txt"}
		if nonXMLFile.IsXML() {
			t.Error("expected IsXML() to return false for non-.xml file")
		}
	})

	t.Run("Exists checks file existence", func(t *testing.T) {
		tmpFile, _ := ioutil.TempFile("", "test")
		defer os.Remove(tmpFile.Name())

		existingFile := &XMLFile{FilesystemPath: tmpFile.Name()}
		if !existingFile.Exists() {
			t.Error("expected Exists() to return true for existing file")
		}

		nonExistentFile := &XMLFile{FilesystemPath: "/this/does/not/exist.xml"}
		if nonExistentFile.Exists() {
			t.Error("expected Exists() to return false for non-existent file")
		}
	})
}

func TestHasExtension(t *testing.T) {
	var tests = []struct {
		path      string
		extension string
		expected  bool
	}{
		{"/path/file.xml", "xml", true},
		{"/path/file.html", "html", true},
		{"/path/file.xml?query", "xml", false}, // Query strings make extension not at end
		{"/path/file.txt", "xml", false},
		{"/path/filexml", "xml", false},
		{"/path/file", "xml", false},
	}

	for _, test := range tests {
		got := hasExtension(test.path, test.extension)
		if got != test.expected {
			t.Errorf("hasExtension(%s, %s) = %v, expected %v", test.path, test.extension, got, test.expected)
		}
	}
}

func TestPathWithoutExt(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/path/file.xml", "/path/file"},
		{"/path/file.html", "/path/file"},
		{"/path/file", "/path/file"},
		{"/file.xml", "/file"},
	}

	for _, test := range tests {
		got := pathWithoutExt(test.input)
		if got != test.expected {
			t.Errorf("pathWithoutExt(%s) = %s, expected %s", test.input, got, test.expected)
		}
	}
}
