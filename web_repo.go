package main

import (
	"errors"
	"net/url"
	"os"
	"path"
	"strings"
)

type WebRepo struct {
	WebRoot string
}

type XMLFile struct {
	URLPath        string
	FilesystemPath string
}

func (x *XMLFile) IsXML() bool {
	return hasExtension(x.FilesystemPath, "xml")
}

func (x *XMLFile) Exists() bool {
	return fileExists(x.FilesystemPath)
}

func NewWebRepo(webRoot string) *WebRepo {
	return &WebRepo{WebRoot: webRoot}
}

func (w *WebRepo) FindXML(requestPath string) (*XMLFile, error) {
	xmlFile, err := w.cleanPath(requestPath)
	if err != nil {
		return nil, err
	}
	if xmlFile.URLPath == "" || xmlFile.URLPath == "/" {
		xmlFile.URLPath = "index.xml"
		xmlFile.FilesystemPath = path.Join(w.WebRoot, "index.xml")
		return xmlFile, nil
	}

	if hasExtension(xmlFile.FilesystemPath, "xml") && fileExists(xmlFile.FilesystemPath) {
		return xmlFile, nil
	}

	if hasExtension(xmlFile.FilesystemPath, "html") {
		x := pathWithoutExt(xmlFile.FilesystemPath) + ".xml"
		if fileExists(x) {
			xmlFile.URLPath = pathWithoutExt(xmlFile.URLPath) + ".xml"
			xmlFile.FilesystemPath = x
		}
		return xmlFile, nil
	}

	if fileExists(xmlFile.FilesystemPath + ".xml") {
		return &XMLFile{
			URLPath:        xmlFile.URLPath + ".xml",
			FilesystemPath: xmlFile.FilesystemPath + ".xml",
		}, nil
	}

	return xmlFile, nil
}

func (w *WebRepo) cleanPath(urlPath string) (*XMLFile, error) {
	urlPath, err := url.PathUnescape(urlPath)
	if err != nil {
		return nil, err
	}
	p := path.Clean(urlPath)
	f, err := prefixWebRoot(w.WebRoot, p)
	if err != nil {
		return nil, err
	}
	return &XMLFile{URLPath: p, FilesystemPath: f}, nil
}

func prefixWebRoot(webRoot string, relative string) (string, error) {
	abs := path.Join(webRoot, relative)
	if strings.Index(abs, webRoot) != 0 {
		return "", errors.New("requested path is not inside webroot")
	}

	return abs, nil
}

func pathWithoutExt(absPath string) string {
	ext := path.Ext(absPath)
	return absPath[0 : len(absPath)-len(ext)]
}

func hasExtension(requestPath, extension string) bool {
	pos := strings.LastIndex(requestPath, extension)
	if pos > -1 && pos == len(requestPath)-len(extension) {
		if requestPath[pos-1] == '.' || requestPath[pos-1] == '?' {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
