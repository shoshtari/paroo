package configs

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type SectionLog struct {
	Environment string `mapstructure:"environment"`
	Level       string `mapstructure:"level"`
}

type SectionHTTPServer struct {
	Address string `mapstructure:"address"`
}

type SectionTelegram struct {
	BaseAddress string        `mapstructure:"base_address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	Token       string        `mapstructure:"token"`
	Proxy       string        `mapstructure:"proxy"`
	ChatID      int           `mapstructure:"chat_id"`
}

type ParooConfig struct {
	Log        SectionLog        `mapstructure:"log"`
	HTTPServer SectionHTTPServer `mapstructure:"http_server"`
	Telegram   SectionTelegram   `mapstructure:"telegram"`
}

func GetConfig(configPaths ...string) (ParooConfig, error) {
	viper.AddConfigPath("/etc/paroo")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
	var ans ParooConfig
	if err := viper.ReadInConfig(); err != nil {
		return ans, errors.Wrap(err, "couldn't read config file")
	}
	if err := viper.Unmarshal(&ans); err != nil {
		return ans, errors.Wrap(err, "couldn't unmarshal config")
	}
	return ans, nil
}
