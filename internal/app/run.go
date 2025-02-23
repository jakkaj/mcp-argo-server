package app

import (
	"fmt"

	"github.com/strowk/foxy-contexts/pkg/app"
	"github.com/strowk/foxy-contexts/pkg/stdio"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	builder := app.NewBuilder().
		WithName("argo-mcp-wrapper").
		WithVersion("0.1.0").
		WithTransport(stdio.NewTransport()).
		WithFxOptions(
			fx.Provide(func() (*zap.Logger, error) {
				return zap.NewDevelopment()
			}),
			fx.Provide(ProvideWorkflowClient),
			fx.Provide(func(logger *zap.Logger) fxeventLogger {
				return fxeventLogger{logger: logger}
				// Removed fx.WithLogger to avoid duplicate fxeventLogger provider
			}),
		)

	// Register tools
	builder = registerLaunchTool(builder)
	builder = registerStatusTool(builder)
	builder = registerResultTool(builder)

	if err := builder.Run(); err != nil {
		fmt.Println("Failed to run argo-mcp-wrapper:", err)
	}
}
