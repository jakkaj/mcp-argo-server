package app

import "github.com/strowk/foxy-contexts/pkg/mcp"

// interface of launch tool
type ArgoToolInterface interface {
	launchHandler(args map[string]interface{}) *mcp.CallToolResult
}
