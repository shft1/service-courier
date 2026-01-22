package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapAdapter - адаптер логгера zap под интерфейс приложения
type ZapAdapter struct {
	logger *zap.Logger
}

// NewZapAdapter - конструктор адаптера zap
func NewZapAdapter() (*ZapAdapter, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.DisableCaller = true
	config.DisableStacktrace = true

	zapLogger, err := config.Build()

	if err != nil {
		return nil, err
	}
	return &ZapAdapter{logger: zapLogger}, nil
}

func (z *ZapAdapter) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Info(msg string, fields ...Field) {
	z.logger.Info(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Error(msg string, fields ...Field) {
	z.logger.Error(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Sync() error {
	return z.logger.Sync()
}

func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	return zapFields
}
