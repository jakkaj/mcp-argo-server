package app

import (
	"encoding/json"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

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

// Wf node outputs to json
// Define output struct types
type ParameterOutput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ArtifactOutput struct {
	Name  string `json:"name"`
	S3Key string `json:"s3Key"`
}

type NodeOutput struct {
	NodeName     string            `json:"nodeName"`
	Parameters   []ParameterOutput `json:"parameters,omitempty"`
	Artifacts    []ArtifactOutput  `json:"artifacts,omitempty"`
	TemplateName string            `json:"templateName,omitempty"`
}

func extractOutputs(wf *v1alpha1.Workflow) (string, error) {
	var nodesOutput []NodeOutput

	for nodeName, node := range wf.Status.Nodes {
		if node.Outputs != nil {
			var params []ParameterOutput
			for _, param := range node.Outputs.Parameters {
				params = append(params, ParameterOutput{
					Name:  param.Name,
					Value: param.Value.String(),
				})
			}

			// Collect artifacts.
			var artifacts []ArtifactOutput
			for _, artifact := range node.Outputs.Artifacts {
				artifacts = append(artifacts, ArtifactOutput{
					Name:  artifact.Name,
					S3Key: artifact.S3.Key,
				})
			}

			nodesOutput = append(nodesOutput, NodeOutput{
				NodeName:     nodeName,
				Parameters:   params,
				Artifacts:    artifacts,
				TemplateName: node.TemplateName,
			})
		}
	}

	result := struct {
		Nodes []NodeOutput `json:"nodes"`
	}{
		Nodes: nodesOutput,
	}

	// Marshal the result into JSON.
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
