package utils

import (
	"go.uber.org/zap"
)

// Logger instance
var Logger *zap.Logger

// InitLogger initializes Zap
func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}
