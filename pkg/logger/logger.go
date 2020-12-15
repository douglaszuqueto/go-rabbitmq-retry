package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logg *zap.Logger

// New New
func New(name string) error {
	config := getConfig()

	defaultFields := zap.Fields(
		zap.String("app", name),
	)

	var err error

	logg, err = config.Build(
		defaultFields,
	)
	if err != nil {
		return err
	}
	defer logg.Sync()

	return nil
}

func getConfig() zap.Config {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     getEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: false,
		DisableCaller:     true,
	}

	return config
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return encoderConfig
}

// Info Info
func Info(msg string, fields ...zap.Field) {
	logg.Info(msg, fields...)
}

// Warning Warning
func Warning(msg string, fields ...zap.Field) {
	logg.Warn(msg, fields...)
}

// Error Error
func Error(msg string, fields ...zap.Field) {
	logg.Error(msg, fields...)
}

// Debug Debug
func Debug(msg string, fields ...zap.Field) {
	logg.Debug(msg, fields...)
}

// Fatal Fatal
func Fatal(msg string, fields ...zap.Field) {
	logg.Fatal(msg, fields...)
}
