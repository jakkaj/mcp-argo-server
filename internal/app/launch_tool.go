package app

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/ghodss/yaml"
	"github.com/strowk/foxy-contexts/pkg/app"
	"github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LaunchTool struct {
	wfClient wfclientset.Interface
	logger   *zap.Logger
}

func NewLaunchTool(wfClient wfclientset.Interface, logger *zap.Logger) *LaunchTool {
	return &LaunchTool{
		wfClient: wfClient,
		logger:   logger,
	}
}

// New extracted handler function to allow reusage from other parts of the application
func (h *LaunchTool) launchHandler(args map[string]interface{}) *mcp.CallToolResult {
	manifestYAML, ok := args["manifest"].(string)
	if !ok || manifestYAML == "" {
		return errorResult("manifest is required and must be a YAML string")
	}
	namespace := "default"
	if nsArg, ok := args["namespace"].(string); ok && nsArg != "" {
		namespace = nsArg
	}

	var wf v1alpha1.Workflow
	if err := yaml.Unmarshal([]byte(manifestYAML), &wf); err != nil {
		h.logger.Error("Invalid workflow YAML", zap.Error(err))
		return errorResult(fmt.Sprintf("Invalid workflow manifest: %v", err))
	}
	if wf.ObjectMeta.Namespace == "" {
		wf.ObjectMeta.Namespace = namespace
	}

	ctx := context.Background()
	createdWf, err := h.wfClient.ArgoprojV1alpha1().Workflows(namespace).Create(ctx, &wf, metav1.CreateOptions{})
	if err != nil {
		h.logger.Error("Failed to create workflow", zap.Error(err))
		return errorResult(fmt.Sprintf("Failed to submit workflow: %v", err))
	}
	h.logger.Info("Workflow submitted", zap.String("name", createdWf.Name), zap.String("namespace", namespace))

	return successResult(fmt.Sprintf("Workflow %q submitted successfully", createdWf.Name))
}

// Change function signature to accept and return a pointer to app.Builder
func registerLaunchTool(builder *app.Builder) *app.Builder {
	return builder.WithTool(func(wfClient wfclientset.Interface, logger *zap.Logger) fxctx.Tool {
		meta := &mcp.Tool{
			Name:        "launch",
			Description: Ptr("Submits a new Argo workflow"),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"manifest":  {"type": "string", "description": "Argo Workflow YAML manifest"},
					"namespace": {"type": "string", "description": "Kubernetes namespace (optional)"},
				},
				Required: []string{"manifest"},
			},
		}
		s := NewLaunchTool(wfClient, logger)
		return fxctx.NewTool(meta, s.launchHandler)
	})
}
