package pkg

import (
	"errors"

	"github.com/shoshtari/paroo/internal/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitializeLogger(config configs.SectionLog) error {

	var loggerConfig zap.Config
	var err error

	switch config.Environment {
	case "production":
		loggerConfig = zap.NewProductionConfig()
	case "development":
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		return errors.New("environment is unknown")
	}

	if err := loggerConfig.Level.UnmarshalText([]byte(config.Level)); err != nil {
		return errors.New("couldn't unmarshal level")
	}
	logger, err = loggerConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

func GetLogger() *zap.Logger {
	return logger
}
