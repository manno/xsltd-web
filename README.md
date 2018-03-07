# XSLTD-CHAOSPAGE

## Dependencies

* xalan-c

## Running locally

```
#!/bin/sh
go get github.com/c4/xsltd-web
go install github.com/c4/xsltd-web

export XALAN=/usr/local/bin/Xalan
export WEBROOT=$HOME/workspace/xsltd-c4/svn/sandbox
xsltd-web
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
