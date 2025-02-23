# Makefile for mcp-argo-server

# Initialize the Go module
init:
	go mod init mcp-argo-server
	go mod tidy

# Build the application
build:
	go build -o bin/mcp-argo-server ./cmd/mcp-argo-server

# Run the application
run: build
	go run ./cmd/mcp-argo-server

# Clean the build
clean:
	rm -rf bin/


clean-k3d:
	k3d cluster delete argoagent-cluster || true
	k3d registry delete k3d-registry.localhost || true
	docker stop k3d-registry.localhost || true
	docker rm k3d-registry.localhost || true

cluster: clean-k3d
	k3d registry create k3d-registry.localhost --port 5000 || true
	k3d cluster create argoagent-cluster \
		--api-port 6443 \
		--servers 1 \
		--agents 1 \
		--registry-use k3d-registry.localhost:5000 \
		--wait \
		-p "2746:32746@server:0"
	./install_argo.sh
	kubectl patch service argo-server -n argo --patch '{"spec": {"type": "NodePort", "ports": [{"name": "web", "port": 2746, "targetPort": 2746, "nodePort": 32746}]}}'
	echo "Argo server is available at https://localhost:2746"
	