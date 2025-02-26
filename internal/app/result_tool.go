package app

import (
	"context"
	"fmt"

	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/strowk/foxy-contexts/pkg/app"
	"github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResultTool struct {
	wfClient wfclientset.Interface
	logger   *zap.Logger
}

func NewResultTool(wfClient wfclientset.Interface, logger *zap.Logger) *ResultTool {
	return &ResultTool{
		wfClient: wfClient,
		logger:   logger,
	}
}

func (h *ResultTool) resultHandler(args map[string]interface{}) *mcp.CallToolResult {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return errorResult("workflow name is required")
	}
	namespace := "default"
	if nsArg, ok := args["namespace"].(string); ok && nsArg != "" {
		namespace = nsArg
	}

	ctx := context.Background()
	wf, err := h.wfClient.ArgoprojV1alpha1().Workflows(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		h.logger.Error("Failed to get workflow for result", zap.Error(err))
		return errorResult(fmt.Sprintf("Failed to get workflow: %v", err))
	}

	phase := string(wf.Status.Phase)
	if phase != "Succeeded" {
		h.logger.Info("Requested results for incomplete or failed workflow",
			zap.String("name", name), zap.String("phase", phase))
		return errorResult(fmt.Sprintf("Workflow %q is not completed (current phase: %s)", name, phase))
	}

	//outputs := wf.Status.Outputs
	var contentItems []interface{}

	// if outputs != nil && len(outputs.Parameters) > 0 {
	// 	for _, param := range outputs.Parameters {
	// 		contentItems = append(contentItems, mcp.TextContent{
	// 			Type: "text",
	// 			Text: fmt.Sprintf("%s: %s", param.Name, param.Value),
	// 		})
	// 	}
	// }

	// if outputs != nil && len(outputs.Artifacts) > 0 {
	// 	for _, art := range outputs.Artifacts {
	// 		info := art.Name
	// 		if art.Path != "" {
	// 			info += fmt.Sprintf(" (path: %s)", art.Path)
	// 		}
	// 		contentItems = append(contentItems, mcp.TextContent{
	// 			Type: "text",
	// 			Text: fmt.Sprintf("Artifact: %s", info),
	// 		})
	// 	}
	// }

	json_outputs, err := extractOutputs(wf)

	if err != nil {
		h.logger.Error("Failed to extract outputs", zap.Error(err))
		return errorResult(fmt.Sprintf("Failed to extract outputs: %v", err))
	}

	contentItems = append(contentItems, mcp.TextContent{
		Type: "text",
		Text: fmt.Sprintf(json_outputs),
	})

	if len(contentItems) == 0 {
		contentItems = append(contentItems, mcp.TextContent{
			Type: "text",
			Text: "No outputs found for workflow.",
		})
	}

	h.logger.Info("Workflow outputs fetched", zap.String("name", name), zap.String("phase", phase))
	return &mcp.CallToolResult{
		IsError: boolPtr(false),
		Content: contentItems,
	}
}

// Updated to accept and return *app.Builder
func registerResultTool(builder *app.Builder) *app.Builder {
	return builder.WithTool(func(wfClient wfclientset.Interface, logger *zap.Logger) fxctx.Tool {
		meta := &mcp.Tool{
			Name:        "result",
			Description: Ptr("Fetches output parameters (and artifacts) from a completed workflow"),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"name":      {"type": "string", "description": "Name of the workflow"},
					"namespace": {"type": "string", "description": "Kubernetes namespace (optional)"},
				},
				Required: []string{"name"},
			},
		}
		s := NewResultTool(wfClient, logger)
		return fxctx.NewTool(meta, s.resultHandler)
	})
}
