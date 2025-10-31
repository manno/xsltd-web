package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandoff(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "handoff-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create xsl subdirectory
	xslDir := filepath.Join(tmpDir, "xsl")
	os.Mkdir(xslDir, 0755)

	// Set up test config
	originalConfig := config
	config = &Config{
		WebRoot: tmpDir,
		Xalan:   "xalan",
		Listen:  "8080",
	}
	defer func() { config = originalConfig }()

	t.Run("handoff with missing stylesheet", func(t *testing.T) {
		// Create XML without stylesheet declaration
		xmlFile := filepath.Join(tmpDir, "no-stylesheet.xml")
		xmlContent := `<?xml version="1.0"?>
<root>Test content</root>`
		ioutil.WriteFile(xmlFile, []byte(xmlContent), 0644)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/no-stylesheet.xml", nil)

		xmlFileObj := &XMLFile{
			URLPath:        "/no-stylesheet.xml",
			FilesystemPath: xmlFile,
		}

		handoff(w, req, xmlFileObj)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		// Should return 500 error because no stylesheet found
		if resp.StatusCode != 500 {
			t.Errorf("expected status 500, got %d", resp.StatusCode)
		}
		if len(body) == 0 {
			t.Error("expected error message in response body")
		}
	})

	t.Run("handoff with stylesheet but xalan not available", func(t *testing.T) {
		// Create XSL file
		xslFile := filepath.Join(xslDir, "test.xsl")
		xslContent := `<?xml version="1.0"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:output method="html"/>
  <xsl:template match="/">
    <html><body>Test</body></html>
  </xsl:template>
</xsl:stylesheet>`
		ioutil.WriteFile(xslFile, []byte(xslContent), 0644)

		// Create XML with stylesheet declaration
		xmlFile := filepath.Join(tmpDir, "with-stylesheet.xml")
		xmlContent := `<?xml version="1.0"?>
<?xml-stylesheet href="/xsl/test.xsl" type="text/xsl"?>
<root>Test content</root>`
		ioutil.WriteFile(xmlFile, []byte(xmlContent), 0644)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/with-stylesheet.xml", nil)

		xmlFileObj := &XMLFile{
			URLPath:        "/with-stylesheet.xml",
			FilesystemPath: xmlFile,
		}

		handoff(w, req, xmlFileObj)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		// Should return 500 error if xalan is not available
		// If xalan is available, it might succeed (status 200)
		if resp.StatusCode != 500 && resp.StatusCode != 200 {
			t.Errorf("expected status 500 or 200, got %d", resp.StatusCode)
		}

		// If status is 500, there should be an error message
		if resp.StatusCode == 500 && len(body) == 0 {
			t.Error("expected error message in response body for 500 status")
		}
	})

	t.Run("handoff with invalid xml file path", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/nonexistent.xml", nil)

		xmlFileObj := &XMLFile{
			URLPath:        "/nonexistent.xml",
			FilesystemPath: filepath.Join(tmpDir, "nonexistent.xml"),
		}

		handoff(w, req, xmlFileObj)

		resp := w.Result()

		// Should return error (500) because file doesn't exist
		if resp.StatusCode != 500 {
			t.Errorf("expected status 500, got %d", resp.StatusCode)
		}
	})
}
