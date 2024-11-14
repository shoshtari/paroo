package internal

import (
	"errors"

	"github.com/shoshtari/paroo/internal/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(config configs.SectionLog) (*zap.Logger, error) {

	var loggerConfig zap.Config
	var err error

	switch config.Environment {
	case "production":
		loggerConfig = zap.NewProductionConfig()
	case "development":
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		return nil, errors.New("environment is unknown")
	}

	if err := loggerConfig.Level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, errors.New("couldn't unmarshal level")
	}
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
