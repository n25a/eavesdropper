package log

import (
	"go.uber.org/zap"
)

var (
	// Logger is the default logger of the application
	Logger *zap.Logger
)

func init() {
	Logger, _ = zap.NewProduction(zap.WithCaller(false))
}
