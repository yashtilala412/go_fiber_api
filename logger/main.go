package logger

import (
	log "github.com/sirupsen/logrus"     // Used for fallback JSON logging
	"go.uber.org/zap"                    // Zap logging package
	"go.uber.org/zap/zapcore"           // Core zap functionalities for advanced configuration
)

// Default encoder configuration used across logging modes.
// Controls how logs are formatted (keys, timestamp format, level formatting, etc.)
var defaultEncoderConfig = zapcore.EncoderConfig{
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

// Default zap logger configuration for server use
var zapServerConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.InfoLevel), // Default log level
	Development:      false,                               // Default to production mode
	Encoding:         "json",                              // Log output format
	EncoderConfig:    defaultEncoderConfig,                // Encoding config
	OutputPaths:      []string{"stdout"},                  // Where to write logs
	ErrorOutputPaths: []string{"stderr"},                  // Where to write errors
}

// NewRootLogger creates and returns a configured zap.Logger instance.
// Supports toggling debug level and development mode formatting.
//
// Parameters:
//   - debug: enables debug-level logging
//   - development: toggles human-friendly colored console output
//
// Returns:
//   - *zap.Logger: configured logger
//   - error: any error during logger construction
func NewRootLogger(debug, development bool) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger

	if debug {
		// Enable debug logging level
		zapServerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		if !development {
			// Debug mode with JSON encoding (production)
			return zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
		}
		// Debug mode with development-friendly colored logs
		zapServerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	if development {
		// Non-debug but development mode: console output with colored levels
		zapServerConfig.Encoding = "console"
		zapServerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	// Default case: production mode with JSON logs and logrus fallback formatter
	log.SetFormatter(&log.JSONFormatter{}) // Set logrus fallback format to JSON
	logger, err = zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	if err != nil {
		panic(err) // Panic on logger initialization failure
	}

	return logger, err
}
