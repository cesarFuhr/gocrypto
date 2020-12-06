package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger a interface to the logger object
type Logger interface {
	Info(...interface{})
}

// NewLogger creates a new logger
func NewLogger() Logger {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}
	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
