package configs

import (
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

type ParooConfig struct {
	Log        SectionLog        `mapstructure:"log"`
	HTTPServer SectionHTTPServer `mapstructure:"http_server"`
}

func GetConfig() (ParooConfig, error) {
	viper.AddConfigPath("/etc/paroo")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	var ans ParooConfig
	if err := viper.ReadInConfig(); err != nil {
		return ans, errors.Wrap(err, "couldn't read config file")
	}
	if err := viper.Unmarshal(&ans); err != nil {
		return ans, errors.Wrap(err, "couldn't unmarshal config")
	}
	return ans, nil
}
