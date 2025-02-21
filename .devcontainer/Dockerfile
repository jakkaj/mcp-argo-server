FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Install required packages
RUN apt-get update && apt-get install -y curl sudo tar python3 python3-pip

# Clean up
RUN rm -rf /var/lib/apt/lists/* /tmp/k9s.tar.gz
