package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// NewRootLogger initializes and returns a new zap.Logger
func NewRootLogger(debug, development bool) (*zap.Logger, error) {
	// Default encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Default logger configuration
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      development,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Adjust settings based on debug & development mode
	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	if development {
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return config.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
}

// GetLogger returns a singleton logger instance
func GetLogger(debug, development bool) *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = NewRootLogger(debug, development)
		if err != nil {
			panic("Failed to initialize logger: " + err.Error()) // Only panics once
		}
	})
	return logger
}
