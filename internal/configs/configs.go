package configs

import (
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type SectionLoggers struct {
	Name  string `mapstructure:"name"`
	Level string `mapstructure:"level"`
}

type SectionLog struct {
	Environment string           `mapstructure:"environment"`
	Level       string           `mapstructure:"level"` // default level, no logger can have a level below this
	Loggers     []SectionLoggers `mapstructure:"loggers"`
}

type SectionHTTPServer struct {
	Address string `mapstructure:"address"`
}

type SectionTelegram struct {
	BaseAddress      string        `mapstructure:"base_address"`
	Timeout          time.Duration `mapstructure:"timeout"`
	Token            string        `mapstructure:"token"`
	Proxy            string        `mapstructure:"proxy"`
	ChatID           int           `mapstructure:"chat_id"`
	GetUpdateTimeout int           `mapstructure:"get_update_timeout"`
}

type SectionWallex struct {
	BaseAddress string        `mapstructure:"base_address"`
	Token       string        `mapstructure:"token"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

type SectionRamzinex struct {
	BaseAddress string        `mapstructure:"base_address"`
	Token       string        `mapstructure:"token"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

type SectionDatabase struct {
	Postgres SectionPostgres `mapstructure:"postgres"`
}

type SectionRedis struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type SectionPostgres struct {
	Host            string        `mapstructure:"host"`
	Port            uint16        `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	User            string        `mapstructure:"user"`
	Pass            string        `mapstructure:"pass"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	ConnMaxTime     time.Duration `mapstructure:"conn_max_time"`
	MaxConn         int32         `mapstructure:"max_conn"`
	MinConn         int32         `mapstructure:"min_conn"`
}

type ParooConfig struct {
	Log        SectionLog        `mapstructure:"log"`
	HTTPServer SectionHTTPServer `mapstructure:"http_server"`
	Telegram   SectionTelegram   `mapstructure:"telegram"`
	Exchange   SectionExchange   `mapstructure:"exchange"`
	Database   SectionDatabase   `mapstructure:"database"`
}

type SectionExchange struct {
	Wallex   SectionWallex   `mapstructure:"wallex"`
	Ramzinex SectionRamzinex `mapstructure:"ramzinex"`
}

func bindEnv(typ reflect.Type, parentpath string) error {
	var basePath string
	if parentpath != "" {
		basePath = parentpath + "."
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldPath := strings.ToLower(basePath + field.Tag.Get("mapstructure"))
		if field.Type.Kind() == reflect.Struct {
			if err := bindEnv(field.Type, fieldPath); err != nil {
				return err
			}
		} else {
			if err := viper.BindEnv(fieldPath); err != nil {
				return err
			}
		}
	}
	return nil

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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := bindEnv(reflect.TypeOf(&ans).Elem(), ""); err != nil {
		return ans, errors.Wrap(err, "couldn't bind env")
	}
	if err := viper.Unmarshal(&ans); err != nil {

		return ans, errors.Wrap(err, "couldn't unmarshal config")
	}
	return ans, nil
}
