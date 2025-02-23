package app

import (
	"fmt"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type fxeventLogger struct {
	logger *zap.Logger
}

func (l fxeventLogger) LogEvent(e fxevent.Event) {
	switch ev := e.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Info("OnStart hook executing", zap.String("callee", ev.FunctionName), zap.String("caller", ev.CallerName))
	case *fxevent.OnStopExecuting:
		l.logger.Info("OnStop hook executing", zap.String("callee", ev.FunctionName), zap.String("caller", ev.CallerName))
	case *fxevent.Provided:
		for _, t := range ev.OutputTypeNames {
			l.logger.Debug("Provided", zap.String("constructor", ev.ConstructorName), zap.String("type", t))
		}
	case *fxevent.Decorated:
		for _, t := range ev.OutputTypeNames {
			l.logger.Debug("Decorated", zap.String("decorator", ev.DecoratorName), zap.String("type", t))
		}
	default:
		l.logger.Debug("Fx event", zap.String("type", fmt.Sprintf("%T", e)))
	}
}
