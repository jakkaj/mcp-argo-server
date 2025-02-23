package app

import (
	"fmt"
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

	//add some more to opts
	opts = append(opts, fx.Invoke(func(tool *LaunchTool) {
		res := tool.launchHandler(map[string]interface{}{
			"manifest":  "test",
			"namespace": "test",
		})
		if res.IsError == nil && !*res.IsError {
			t.Error("Expected no error, got:", res.Content)
		}
	}))

	app := fxtest.New(t, opts...)

	defer app.RequireStart().RequireStop()

}
