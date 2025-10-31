package main

import (
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	var tests = []struct {
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{"TEST_KEY", "custom_value", "default", "custom_value"},
		{"UNSET_KEY", "", "default", "default"},
	}

	for _, test := range tests {
		if test.envValue != "" {
			os.Setenv(test.envKey, test.envValue)
			defer os.Unsetenv(test.envKey)
		}

		got := fetch(test.envKey, test.defaultValue)
		if got != test.expected {
			t.Errorf("fetch(%s, %s) = %s, expected %s", test.envKey, test.defaultValue, got, test.expected)
		}
	}
}

func TestConfigSetup(t *testing.T) {
	// Save original env vars
	originalListen := os.Getenv("LISTEN")
	originalWebRoot := os.Getenv("WEBROOT")
	originalXalan := os.Getenv("XALAN")

	// Clear env vars
	os.Unsetenv("LISTEN")
	os.Unsetenv("WEBROOT")
	os.Unsetenv("XALAN")

	// Test with defaults
	cfg := &Config{}
	cfg.Setup()

	if cfg.Listen != "8080" {
		t.Errorf("Expected default Listen to be '8080', got '%s'", cfg.Listen)
	}
	if cfg.WebRoot != "/home/www/koeln.ccc.de/sandbox" {
		t.Errorf("Expected default WebRoot, got '%s'", cfg.WebRoot)
	}
	if cfg.Xalan != "xalan" {
		t.Errorf("Expected default Xalan to be 'xalan', got '%s'", cfg.Xalan)
	}

	// Test with custom env vars
	os.Setenv("LISTEN", "localhost:9000")
	os.Setenv("WEBROOT", "/custom/path")
	os.Setenv("XALAN", "/usr/bin/custom-xalan")

	cfg = &Config{}
	cfg.Setup()

	if cfg.Listen != "localhost:9000" {
		t.Errorf("Expected Listen to be 'localhost:9000', got '%s'", cfg.Listen)
	}
	if cfg.WebRoot != "/custom/path" {
		t.Errorf("Expected WebRoot to be '/custom/path', got '%s'", cfg.WebRoot)
	}
	if cfg.Xalan != "/usr/bin/custom-xalan" {
		t.Errorf("Expected Xalan to be '/usr/bin/custom-xalan', got '%s'", cfg.Xalan)
	}

	// Restore original env vars
	if originalListen != "" {
		os.Setenv("LISTEN", originalListen)
	} else {
		os.Unsetenv("LISTEN")
	}
	if originalWebRoot != "" {
		os.Setenv("WEBROOT", originalWebRoot)
	} else {
		os.Unsetenv("WEBROOT")
	}
	if originalXalan != "" {
		os.Setenv("XALAN", originalXalan)
	} else {
		os.Unsetenv("XALAN")
	}
}

func TestConfigValidate(t *testing.T) {
	// Test with non-existent WebRoot
	cfg := &Config{
		Listen:  "8080",
		WebRoot: "/this/path/does/not/exist",
		Xalan:   "xalan",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for non-existent WebRoot, got nil")
	}

	// Test with non-existent xalan binary
	tmpDir := os.TempDir()
	cfg = &Config{
		Listen:  "8080",
		WebRoot: tmpDir,
		Xalan:   "this-binary-does-not-exist-12345",
	}

	err = cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for non-existent xalan binary, got nil")
	}
}
