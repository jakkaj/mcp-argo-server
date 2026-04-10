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

type DeleteWorkflowTool struct {
	wfClient wfclientset.Interface
	logger   *zap.Logger
}

func NewDeleteWorkflowTool(wfClient wfclientset.Interface, logger *zap.Logger) *DeleteWorkflowTool {
	return &DeleteWorkflowTool{
		wfClient: wfClient,
		logger:   logger,
	}
}

func (h *DeleteWorkflowTool) deleteWorkflowHandler(args map[string]interface{}) *mcp.CallToolResult {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return errorResult("workflow name is required")
	}

	namespace := "argo"
	if nsArg, ok := args["namespace"].(string); ok && nsArg != "" {
		namespace = nsArg
	}

	if err := h.wfClient.ArgoprojV1alpha1().Workflows(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		h.logger.Error("Failed to delete workflow", zap.String("name", name), zap.String("namespace", namespace), zap.Error(err))
		return errorResult(fmt.Sprintf("Failed to delete workflow %q: %v", name, err))
	}

	h.logger.Info("Workflow deleted", zap.String("name", name), zap.String("namespace", namespace))
	return successResult(fmt.Sprintf("Workflow %q deleted from namespace %q", name, namespace))
}

func registerDeleteWorkflowTool(builder *app.Builder) *app.Builder {
	return builder.WithTool(func(wfClient wfclientset.Interface, logger *zap.Logger) fxctx.Tool {
		meta := &mcp.Tool{
			Name:        "delete_workflow",
			Description: Ptr("Deletes a workflow by name"),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"name":      {"type": "string", "description": "Name of the workflow"},
					"namespace": {"type": "string", "description": "Kubernetes namespace (optional)"},
				},
				Required: []string{"name"},
			},
		}

		s := NewDeleteWorkflowTool(wfClient, logger)
		return fxctx.NewTool(meta, s.deleteWorkflowHandler)
	})
}

type DeleteWorkflowTemplateTool struct {
	wfClient wfclientset.Interface
	logger   *zap.Logger
}

func NewDeleteWorkflowTemplateTool(wfClient wfclientset.Interface, logger *zap.Logger) *DeleteWorkflowTemplateTool {
	return &DeleteWorkflowTemplateTool{
		wfClient: wfClient,
		logger:   logger,
	}
}

func (h *DeleteWorkflowTemplateTool) deleteWorkflowTemplateHandler(args map[string]interface{}) *mcp.CallToolResult {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return errorResult("workflow template name is required")
	}

	namespace := "argo"
	if nsArg, ok := args["namespace"].(string); ok && nsArg != "" {
		namespace = nsArg
	}

	if err := h.wfClient.ArgoprojV1alpha1().WorkflowTemplates(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		h.logger.Error("Failed to delete workflow template", zap.String("name", name), zap.String("namespace", namespace), zap.Error(err))
		return errorResult(fmt.Sprintf("Failed to delete workflow template %q: %v", name, err))
	}

	h.logger.Info("Workflow template deleted", zap.String("name", name), zap.String("namespace", namespace))
	return successResult(fmt.Sprintf("WorkflowTemplate %q deleted from namespace %q", name, namespace))
}

func registerDeleteWorkflowTemplateTool(builder *app.Builder) *app.Builder {
	return builder.WithTool(func(wfClient wfclientset.Interface, logger *zap.Logger) fxctx.Tool {
		meta := &mcp.Tool{
			Name:        "delete_workflow_template",
			Description: Ptr("Deletes a workflow template by name"),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"name":      {"type": "string", "description": "Name of the workflow template"},
					"namespace": {"type": "string", "description": "Kubernetes namespace (optional)"},
				},
				Required: []string{"name"},
			},
		}

		s := NewDeleteWorkflowTemplateTool(wfClient, logger)
		return fxctx.NewTool(meta, s.deleteWorkflowTemplateHandler)
	})
}
