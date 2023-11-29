package util

import (
	"go.uber.org/zap"
)

func GetLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err) // Handle the error according to your needs
	}
	defer logger.Sync() // Flushes buffer, if any
	return logger
}
