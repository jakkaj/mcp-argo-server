package app

import (
	"context"
	"fmt"
	"time"

	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/strowk/foxy-contexts/pkg/app"
	"github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusTool struct {
	wfClient wfclientset.Interface
	logger   *zap.Logger
}

func NewStatusTool(wfClient wfclientset.Interface, logger *zap.Logger) *StatusTool {
	return &StatusTool{
		wfClient: wfClient,
		logger:   logger,
	}
}

func (h *StatusTool) statusHandler(args map[string]interface{}) *mcp.CallToolResult {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return errorResult("workflow name is required")
	}
	namespace := "argo"
	if nsArg, ok := args["namespace"].(string); ok && nsArg != "" {
		namespace = nsArg
	}

	ctx := context.Background()
	wf, err := h.wfClient.ArgoprojV1alpha1().Workflows(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		h.logger.Error("Failed to get workflow status", zap.Error(err))
		return errorResult(fmt.Sprintf("Failed to get workflow status: %v", err))
	}

	phase := string(wf.Status.Phase)
	statusMsg := wf.Status.Message
	if !wf.Status.FinishedAt.IsZero() {
		finishedTime := wf.Status.FinishedAt.Time.Format(time.RFC3339)
		if statusMsg != "" {
			statusMsg = fmt.Sprintf("%s (finished at %s)", statusMsg, finishedTime)
		} else {
			statusMsg = fmt.Sprintf("Finished at %s", finishedTime)
		}
	}

	var outputText string
	if statusMsg != "" {
		outputText = fmt.Sprintf("Workflow %q status: %s - %s", name, phase, statusMsg)
	} else {
		outputText = fmt.Sprintf("Workflow %q status: %s", name, phase)
	}
	h.logger.Info("Workflow status retrieved", zap.String("name", name), zap.String("phase", phase))
	res := successResult(outputText)
	res.Content = append(res.Content, mcp.TextContent{Type: "text", Text: name})
	res.Content = append(res.Content, mcp.TextContent{Type: "text", Text: phase})
	res.Content = append(res.Content, mcp.TextContent{Type: "text", Text: statusMsg})
	return res
}

// Updated to accept and return *app.Builder
func registerStatusTool(builder *app.Builder) *app.Builder {
	return builder.WithTool(func(wfClient wfclientset.Interface, logger *zap.Logger) fxctx.Tool {
		meta := &mcp.Tool{
			Name:        "status",
			Description: Ptr("Gets the status of a workflow by name"),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"name":      {"type": "string", "description": "Name of the workflow"},
					"namespace": {"type": "string", "description": "Kubernetes namespace (optional)"},
				},
				Required: []string{"name"},
			},
		}

		s := NewStatusTool(wfClient, logger)
		return fxctx.NewTool(meta, s.statusHandler)
	})
}
