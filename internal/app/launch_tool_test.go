package app

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSample(t *testing.T) {
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
		fmt.Println(res.Content...)
	}))

	app := fxtest.New(t, opts...)

	defer app.RequireStart().RequireStop()

}
