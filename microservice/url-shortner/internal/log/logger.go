package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the logger based on the config
func InitLogger(logLevel, mode string) {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		panic("Failed to initialize logger = " + err.Error())
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(lvl),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: mode == "release",
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "timestamp",
			LevelKey:         "level",
			NameKey:          "logger",
			CallerKey:        "caller",
			MessageKey:       "message",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalColorLevelEncoder,
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " | ",
		},
		OutputPaths:      []string{"stdout", "/tmp/logs"},
		ErrorOutputPaths: []string{"stderr"},
	}
	var err error
	Logger, err = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic("Failed to initialize logger = " + err.Error())
	}
}

// Info logs an info message with additional context fields
func Info(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

// Error logs an error message with additional context fields
func Error(message string, fields ...zap.Field) {
	Logger.Error(message, fields...)
}

// Debug logs the debug message with additional context fields
func Debug(message string, fields ...zap.Field) {
	Logger.Debug(message, fields...)
}

// Fatal logs a fatal methods with additional context fields
func Fatal(message string, fields ...zap.Field) {
	Logger.Fatal(message, fields...)
}

// Warn logs a warning message with additional context fields
func Warn(message string, fields ...zap.Field) {
	Logger.Warn(message, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Logger.Sync()
}
