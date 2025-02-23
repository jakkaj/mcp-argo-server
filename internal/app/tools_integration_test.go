package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestToolsIntegration(t *testing.T) {
	mt := &SpecialTesting{T: t}

	// Now you can call your extended method.
	opts := mt.Opts()

	if err := fx.ValidateApp(opts...); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Validation passed")
	}

	//load the file argo-hello-world.yaml in to a string
	manifest, err := os.ReadFile("../../kube/argo-hello-world.yaml")
	if err != nil {
		t.Error(err)
	}
	_ = manifest
	//add some more to opts
	opts = append(opts, fx.Invoke(func(tool *LaunchTool) {
		res := tool.launchHandler(map[string]interface{}{
			"manifest":  "test",
			"namespace": "test",
		})
		if res.IsError == nil && !*res.IsError {
			t.Error("Expected error, got:", res.Content)
		}
	}))

	opts = append(opts, fx.Invoke(func(tool *LaunchTool, statusTool *StatusTool, resultTool *ResultTool, wfClient wfclientset.Interface) {
		res := tool.launchHandler(map[string]interface{}{
			"manifest":  string(manifest),
			"namespace": "argo",
		})
		if res.IsError != nil && *res.IsError {
			t.Error("Expected no error, got:", res.Content)
		}

		var wfName string

		for _, item := range res.Content {
			tContent, ok := item.(mcp.TextContent)
			if !ok {
				t.Error("Expected TextContent, got:", res.Content)
			}
			if tContent.Type == "name" {
				wfName = tContent.Text
			}
			fmt.Printf("%s: %s\n", tContent.Type, tContent.Text)
		}
		if wfName == "" {
			t.Error("Expected workflow name, got:", res.Content)
		}

		statusWrong := statusTool.statusHandler(map[string]interface{}{
			"name":      "sdlkfjsdjklsdfjkl",
			"namespace": "argo",
		})

		if statusWrong.IsError == nil && !*statusWrong.IsError {
			t.Error("Expected error, got:", statusWrong.Content)
		}

		status := statusTool.statusHandler(map[string]interface{}{
			"name":      wfName,
			"namespace": "argo",
		})
		if status.IsError != nil && *status.IsError {
			t.Error("Expected no error, got:", status.Content)
		}
		for _, item := range status.Content {
			tContent, ok := item.(mcp.TextContent)
			if !ok {
				t.Error("Expected TextContent, got:", status.Content)
			}
			fmt.Printf("%s: %s\n", tContent.Type, tContent.Text)
		}

		watch, err := wfClient.ArgoprojV1alpha1().Workflows("argo").Watch(
			context.Background(),
			metav1.ListOptions{
				FieldSelector: fmt.Sprintf("metadata.name=%s", wfName),
			},
		)
		if err != nil {
			t.Error(err)
		}
		defer watch.Stop()
		// Wait for the workflow to complete
		for event := range watch.ResultChan() {
			wf, ok := event.Object.(*wfv1.Workflow)
			if !ok {
				log.Printf("Unexpected type in watch event")
				continue
			}
			if wf.Status.Phase == wfv1.WorkflowSucceeded ||
				wf.Status.Phase == wfv1.WorkflowFailed ||
				wf.Status.Phase == wfv1.WorkflowError {
				break
			}

		}
		result := resultTool.resultHandler(map[string]interface{}{
			"name":      wfName,
			"namespace": "argo",
		})

		if result.IsError != nil && *result.IsError {
			t.Error("Expected no error, got:", result.Content)
		}
		for _, item := range result.Content {
			tContent, ok := item.(mcp.TextContent)
			if !ok {
				t.Error("Expected TextContent, got:", result.Content)
			}
			fmt.Printf("%s: %s\n", tContent.Type, tContent.Text)
		}
	}))

	app := fxtest.New(t, opts...)

	defer app.RequireStart().RequireStop()

}
