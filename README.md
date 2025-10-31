# XSLTD-CHAOSPAGE

A Go HTTP server that serves XML files transformed via XSLT stylesheets using the Xalan-C processor.

## Dependencies

* **Go 1.16+**
* **xalan-c** - XSLT processor (install via package manager)

## Installation

### From Source

```bash
git clone https://github.com/c4/xsltd-web
cd xsltd-web
go build
```

### Via Go Install

```bash
go install github.com/c4/xsltd-web@latest
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

# Run tests with coverage
go test -cover

# Run specific test
go test -run TestFindXML -v
```

Current test coverage: **84.2%**

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
