package app

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func toolsModule() fx.Option {
	m := fx.Module("tests",
		fx.Provide(NewLaunchTool),
		fx.Provide(NewStatusTool),
		fx.Provide(NewResultTool),
	)

	return m
}

func argoModule() fx.Option {
	m := fx.Module("argo",
		fx.Provide(ProvideWorkflowClient))

	return m
}

func loggerModule() fx.Option {
	m := fx.Module("logger",
		fx.Provide(func() (*zap.Logger, error) {
			return zap.NewDevelopment()
		}),

		fx.Provide(func(logger *zap.Logger) fxeventLogger {
			return fxeventLogger{logger: logger}
			// Removed fx.WithLogger to avoid duplicate fxeventLogger provider
		}),
	)

	return m
}
