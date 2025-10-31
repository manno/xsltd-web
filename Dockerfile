# Use Debian slim as base - minimal but includes xalan
FROM debian:bookworm-slim

# Install xalan for XSLT processing
RUN apt-get update && \
    apt-get install -y --no-install-recommends xalan && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1000 xsltd && \
    useradd -r -u 1000 -g xsltd xsltd

# Copy the binary from goreleaser build context
COPY xsltd-web /usr/local/bin/xsltd-web

# Set ownership
RUN chown xsltd:xsltd /usr/local/bin/xsltd-web

# Switch to non-root user
USER xsltd

# Set default environment variables
ENV LISTEN=0.0.0.0:8080 \
    WEBROOT=/srv/www \
    XALAN=/usr/bin/Xalan

# Expose default port
EXPOSE 8080

# Run the server
ENTRYPOINT ["/usr/local/bin/xsltd-web"]
