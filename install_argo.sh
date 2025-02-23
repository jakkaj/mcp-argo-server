#!/bin/bash
set -e

# Create namespace for argo if it doesn't exist
kubectl get namespace argo >/dev/null 2>&1 || kubectl create namespace argo

# Install Argo Workflows using the provided minimal quick-start manifest
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.5.4/quick-start-minimal.yaml

# Create ClusterRole and ClusterRoleBinding for argo-server
echo "Creating RBAC rules for argo-server..."
kubectl create clusterrole argo-server-role --verb=get,list,watch,create,update,patch,delete --resource=workflows.argoproj.io,workflowtemplates.argoproj.io,cronworkflows.argoproj.io,clusterworkflowtemplates.argoproj.io --dry-run=client -o yaml | kubectl apply -f -

kubectl create clusterrolebinding argo-server-rb --clusterrole=argo-server-role --serviceaccount=argo:argo-server --dry-run=client -o yaml | kubectl apply -f -

# Wait for Argo components to be up
echo "Waiting for Argo components to start..."
kubectl wait --for=condition=available --timeout=300s deployment/workflow-controller -n argo
kubectl wait --for=condition=available --timeout=300s deployment/argo-server -n argo

# Check if Argo CLI is installed and working, if not install/reinstall it
if ! command -v argo &> /dev/null || ! argo version &> /dev/null; then
    echo "Argo CLI not found or invalid. Installing..."
    # Remove old binary if exists
    if command -v argo &> /dev/null; then
        sudo rm -f "$(which argo)"
    fi
    
    # Set the desired version
    VERSION=$(curl --silent "https://api.github.com/repos/argoproj/argo-workflows/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    
    # Download the binary
    curl -sSL -o argo-linux-amd64.gz https://github.com/argoproj/argo-workflows/releases/download/$VERSION/argo-linux-amd64.gz
    
    # Unzip the binary
    gunzip argo-linux-amd64.gz
    
    # Verify the downloaded file is a valid ELF binary if 'file' command exists
    FILE_CMD=$(command -v file || echo "")
    if [ -n "$FILE_CMD" ]; then
        if ! $FILE_CMD argo-linux-amd64 | grep -q ELF; then
            echo "Downloaded file is not a valid binary."; exit 1
        fi
    else
        echo "Skipping ELF binary verification as 'file' command is not available."
    fi
    
    # Make the binary executable
    chmod +x argo-linux-amd64
    
    # Move the binary to a directory in your PATH
    sudo mv argo-linux-amd64 /usr/local/bin/argo
    
    # Verify the installation
    argo version
fi

# Configure Argo server for insecure access
echo "Configuring Argo Server for insecure access..."
kubectl patch deployment \
    argo-server \
    --namespace argo \
    --type='json' \
    -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/args", "value": ["server", "--auth-mode=server", "--secure=false"]}]'

# Wait for the patch to take effect
echo "Waiting for Argo Server to be ready after configuration change..."
kubectl wait --for=condition=available --timeout=300s deployment/argo-server -n argo
