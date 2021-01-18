package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger a interface to the logger object
type Logger interface {
	Info(string, ...zap.Field)
}

// NewLogger creates a new logger
func NewLogger() Logger {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}
	logConfig.EncoderConfig.EncodeDuration = func(d time.Duration, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(d.String())
	}
	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
