package app

import "github.com/strowk/foxy-contexts/pkg/mcp"

// Ptr returns a pointer to the given string.
func Ptr(s string) *string {
	return &s
}

// boolPtr returns a pointer to a bool.
func boolPtr(b bool) *bool {
	return &b
}

// errorResult returns an error tool result with the provided message.
func errorResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: boolPtr(true),
		Content: []interface{}{
			mcp.TextContent{Type: "text", Text: message},
		},
	}
}

// successResult returns a success tool result with the provided message.
func successResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: boolPtr(false),
		Content: []interface{}{
			mcp.TextContent{Type: "text", Text: message},
		},
	}
}
