package cfg

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func CreateConfigFromViper(v *viper.Viper) (Config, error) {
	var conf Config
	return conf, v.UnmarshalExact(&conf)
}

func NewConfig() (Config, error) {
	if err := ParseAllFlags(); err != nil {
		return Config{}, fmt.Errorf("parse:flags:all:error:%v", err)
	}

	conf, err := CreateConfigFromViper(viper.GetViper())
	if err != nil {
		return Config{}, fmt.Errorf("create:config:from:viper:error:%v", err)
	}

	return conf, nil
}

type Config struct {
	ServeAddr   string `mapstructure:"serve-addr"`
	ServiceName string `mapstructure:"serve-name"`
	DebugMode   bool   `mapstructure:"debug"`

	Postgre `mapstructure:",squash"`
	Rabbit  `mapstructure:",squash"`
}

type Postgre struct {
	URL string `mapstructure:"postgre-url"`
}

type Rabbit struct {
	DSN       string        `mapstructure:"rabbitmq-dsn"`
	Exchange  string        `mapstructure:"rabbit-exchange"`
	Heartbeat time.Duration `mapstructure:"rabbit-heartbeat"`
}
