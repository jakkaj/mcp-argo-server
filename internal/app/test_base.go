package app

import (
	"testing"

	"go.uber.org/fx"
)

type SpecialTesting struct {
	*testing.T
}

func (mt *SpecialTesting) Opts() []fx.Option {
	return []fx.Option{
		toolsModule(),
		argoModule(),
		loggerModule(),
	}
}
