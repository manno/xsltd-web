package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestViewHandler(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "webroot-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up test config
	originalConfig := config
	config = &Config{
		WebRoot: tmpDir,
		Xalan:   "xalan",
		Listen:  "8080",
	}
	defer func() { config = originalConfig }()

	// Create test files
	indexXML := filepath.Join(tmpDir, "index.xml")
	pageXML := filepath.Join(tmpDir, "page.xml")
	textFile := filepath.Join(tmpDir, "test.txt")

	ioutil.WriteFile(indexXML, []byte("<?xml version=\"1.0\"?>\n<root/>"), 0644)
	ioutil.WriteFile(pageXML, []byte("<?xml version=\"1.0\"?>\n<page/>"), 0644)
	ioutil.WriteFile(textFile, []byte("plain text content"), 0644)

	t.Run("serves non-xml files directly", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test.txt", nil)
		w := httptest.NewRecorder()

		viewHandler(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %d", resp.StatusCode)
		}
		if string(body) != "plain text content" {
			t.Errorf("expected 'plain text content', got '%s'", string(body))
		}
	})

	t.Run("returns 404 for non-existent files", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/nonexistent.xml", nil)
		w := httptest.NewRecorder()

		viewHandler(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("handles root path request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		viewHandler(w, req)

		// Note: This will fail XSLT transformation since no stylesheet is defined
		// but the handler should not panic
		resp := w.Result()

		// Should return either 500 (no stylesheet) or 200 (if xalan exists and works)
		if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 500 or 200, got %d", resp.StatusCode)
		}
	})
}

func TestViewHandlerWithXMLStylesheet(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "webroot-xsl-test")
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

	// Create a simple XSL stylesheet
	xslFile := filepath.Join(xslDir, "test.xsl")
	xslContent := `<?xml version="1.0"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:output method="html"/>
  <xsl:template match="/">
    <html><body>Transformed</body></html>
  </xsl:template>
</xsl:stylesheet>`
	ioutil.WriteFile(xslFile, []byte(xslContent), 0644)

	// Create XML file with stylesheet reference
	xmlFile := filepath.Join(tmpDir, "test.xml")
	xmlContent := `<?xml version="1.0"?>
<?xml-stylesheet href="/xsl/test.xsl" type="text/xsl"?>
<root>Test content</root>`
	ioutil.WriteFile(xmlFile, []byte(xmlContent), 0644)

	t.Run("xml with stylesheet processes through handoff", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test.xml", nil)
		w := httptest.NewRecorder()

		viewHandler(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		// This will fail if xalan is not installed, which is expected
		// We're mainly testing that the handler doesn't panic
		if resp.StatusCode == http.StatusInternalServerError {
			// Expected if xalan is not available
			if len(body) == 0 {
				t.Error("expected error message in response body")
			}
		}
		// If xalan is available, status should be 200 or 500 depending on transformation
	})
}
