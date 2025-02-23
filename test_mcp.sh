#!/bin/bash

set -e

# This script tests if the MCP service is working by querying for a list of available tools

# Build the binary if it doesn't exist
if [ ! -f bin/mcp-argo-server ]; then
    echo "Building binary..."
    make build
fi

echo "Sending JSON-RPC request to list available tools..."

# Send a JSON-RPC request for 'list-tools' to the MCP service
response=$(echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./bin/mcp-argo-server)

echo "Response:"
echo "$response"
