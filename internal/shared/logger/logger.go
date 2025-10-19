package logger

import "go.uber.org/zap"

const (
	envLocal = "local"
	envProd  = "prod"
)

func New(logLevel string) *zap.SugaredLogger {
	var logger *zap.Logger
	switch logLevel {
	case envLocal:
		logger, _ = zap.NewDevelopment()
	case envProd:
		logger, _ = zap.NewProduction()
	default:
		logger, _ = zap.NewDevelopment()
	}

	sugar := logger.Sugar()
	return sugar

}
