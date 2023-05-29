package util

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}
