package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/strowk/foxy-contexts/pkg/mcp"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestLaunchTool(t *testing.T) {
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

	opts = append(opts, fx.Invoke(func(tool *LaunchTool) {
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

	}))

	app := fxtest.New(t, opts...)

	defer app.RequireStart().RequireStop()

}
