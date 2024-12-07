package pkg

import (
	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggers map[string]*zap.Logger

func InitializeLogger(config configs.SectionLog) error {

	var loggerConfig zap.Config
	var err error

	loggers = make(map[string]*zap.Logger)
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
	loggers["default"], err = loggerConfig.Build()
	if err != nil {
		return err
	}

	for _, loggerConfig := range config.Loggers {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(loggerConfig.Level)); err != nil {
			return errors.WithStack(err)
		}

		loggers[loggerConfig.Name] = loggers["default"].
			With(zap.String("logger", loggerConfig.Name)).
			WithOptions(zap.IncreaseLevel(level))

	}

	return nil
}

func GetLogger(loggernames ...string) *zap.Logger {
	var loggername string
	if len(loggernames) == 0 {
		loggername = "default"
	} else {
		loggername = loggernames[0]
	}

	if loggers == nil {
		loggers = make(map[string]*zap.Logger)

	}
	if logger, exists := loggers[loggername]; exists {
		return logger
	}

	var err error
	loggers[loggername], err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return loggers[loggername]
}
