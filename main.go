// main.go
//
// Example implementation of an MCP-compliant CLI wrapper for Argo Workflows.
// It uses Foxy Contexts to handle JSON-RPC over STDIN/STDOUT, providing
// three tools: "launch", "status", and "result". Each tool calls Argo
// Workflows via client-go.
//
// This is a minimal example to illustrate the core concepts from the spec.
// Customize logging, error handling, and artifact retrieval as needed.
//
// To build and run:
//
//   go mod init mcp-argo-server
//   go mod tidy
//   go run main.go
//
// Ensure you have Kubernetes/Argo credentials set up correctly (KUBECONFIG
// or in-cluster config) so the wrapper can create/read Workflow resources.

package main

import "mcp-argo-server/internal/app"

func main() {
	app.Run()
}
