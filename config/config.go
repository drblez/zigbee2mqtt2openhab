package config

import (
	"os"
	"path/filepath"

	"github.com/joomcode/errorx"
	"github.com/spf13/viper"
)

var (
	Errors       = errorx.NewNamespace("config")
	CommonErrors = Errors.NewType("common")
)

type Config struct {
	v *viper.Viper
}

func New() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("zigbee2mqtt2openhab.yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("~/.zigbee2mqtt2openhab")
	v.AddConfigPath("/etc/zigbee2mqtt2openhab")
	v.AddConfigPath(filepath.Dir(os.Args[0]))
	v.SetDefault(loggerLevel, defaultLoggerLevel)
	v.SetDefault(z2mTopic, defaultZ2MTopic)
	v.SetDefault(openHABTopic, defaultOpenHABTopic)
	if err := v.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return &Config{v: v}, nil
}

const (
	mqZ2MAddress     = "mq.z2m.address"
	mqOpenHABAddress = "mq.openhab.address"

	loggerLevel        = "log.level"
	defaultLoggerLevel = "debug"

	z2mTopic        = "mq.z2m.topic"
	defaultZ2MTopic = "zigbee2mqtt/oh"

	openHABTopic        = "mq.openhab.topic"
	defaultOpenHABTopic = "openhab"
)

func (cfg *Config) MQZ2MAddress() string {
	return cfg.v.GetString(mqZ2MAddress)
}

func (cfg *Config) MQOpenHABAddress() string {
	return cfg.v.GetString(mqOpenHABAddress)
}

func (cfg *Config) MQZ2MTopic() string {
	return cfg.v.GetString(z2mTopic)
}

func (cfg *Config) MQOpenHABTopic() string {
	return cfg.v.GetString(openHABTopic)
}

func (cfg *Config) LoggerLevel() string {
	return cfg.v.GetString(loggerLevel)
}
