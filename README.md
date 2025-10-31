# XSLTD-CHAOSPAGE

A Go HTTP server that serves XML files transformed via XSLT stylesheets using the Xalan-C processor.

## Dependencies

* **Go 1.16+**
* **xalan-c** - XSLT processor (install via package manager)

## Installation

### Pre-built Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/manno/xsltd-web/releases):

```bash
# Linux (amd64)
wget https://github.com/manno/xsltd-web/releases/latest/download/xsltd-web_Linux_x86_64.tar.gz
tar xzf xsltd-web_Linux_x86_64.tar.gz
chmod +x xsltd-web

# macOS (arm64)
wget https://github.com/manno/xsltd-web/releases/latest/download/xsltd-web_Darwin_arm64.tar.gz
tar xzf xsltd-web_Darwin_arm64.tar.gz
chmod +x xsltd-web
```

### Container Image

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/manno/xsltd-web:latest

# Run with volume mount for your content
docker run -d \
  -p 8080:8080 \
  -v /path/to/your/content:/srv/www \
  -e WEBROOT=/srv/www \
  ghcr.io/manno/xsltd-web:latest
```

The container image is based on Debian slim (~80MB) and includes xalan.

### From Source

```bash
git clone https://github.com/manno/xsltd-web
cd xsltd-web
go build
```

### Via Go Install

```bash
go install github.com/manno/xsltd-web@latest
```

## Running Locally

### Basic Usage

```bash
# Set environment variables
export XALAN=/usr/local/bin/Xalan
export WEBROOT=$HOME/workspace/xsltd-c4/svn/sandbox
export LISTEN=localhost:8080

# Run the server
./xsltd-web
```

### Content Development

**Content changes are reflected immediately** - no server restart needed!

The server reads XML and XSL files from disk on every request (no caching), so:
- Edit your XML files in WEBROOT
- Edit your XSL stylesheets
- Refresh your browser to see changes

This makes content development very fast and iterative.

## Configuration

The server is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN` | `8080` | Server bind address (e.g., `localhost:8080`) |
| `WEBROOT` | `/home/www/koeln.ccc.de/sandbox` | Document root directory |
| `XALAN` | `xalan` | Path to Xalan-C binary |

## Testing

```bash
# Run all tests
go test

# Run specific test
go test -run TestFindXML -v
```

## Run via systemd

```
[Unit]
Description=xsltd - serving chaospage XML with XSL
After=network.target

[Service]
Type=simple
Environment=WEBROOT=/srv/www/chaospages
Environment=LISTEN=localhost:8123
User=www-data
ExecStart=/usr/local/bin/xsltd-web
Restart=always
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=xsltd-chaospage

[Install]
WantedBy=default.target
```

## Releases

Releases are automated via GitHub Actions and [GoReleaser](https://goreleaser.com/):

1. Tag a new version: `git tag -a v1.0.0 -m "Release v1.0.0"`
2. Push the tag: `git push origin v1.0.0`
3. GitHub Actions will automatically:
   - Run tests
   - Build binaries for Linux, macOS, Windows (amd64 & arm64)
   - Create container images for linux/amd64 and linux/arm64
   - Push images to `ghcr.io/manno/xsltd-web`
   - Create a GitHub release with binaries and changelog

### Testing Releases Locally

Install GoReleaser and test the release process locally:

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser
# or
go install github.com/goreleaser/goreleaser/v2@latest

# Validate configuration
goreleaser check

# Build binaries only (fast, no Docker)
goreleaser build --snapshot --clean

# Full release dry-run with Docker images
goreleaser release --snapshot --clean --skip=publish

# Full dry-run without Docker (faster)
goreleaser release --snapshot --clean --skip=publish --skip=docker
```

Artifacts will be in the `./dist/` directory.
