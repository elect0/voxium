package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env string) (*zap.Logger, error) {
	var zapLogger *zap.Logger
	var err error

	switch env {
	case "development":
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = config.Build()

	case "production":
		config := zap.NewProductionConfig()
		config.DisableStacktrace = true
		zapLogger, err = config.Build()
	default:
		return nil, fmt.Errorf("Unknown environment: %s", env)
	}

	if err != nil {
		return nil, err
	}

	return zapLogger, err
}
