package adapters

import (
	"__template__/internal/ports"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

var _ ports.Config = &Config{}

// New returns a new config instance.
func NewConfig() (*Config, error) {
	config := viper.New()
	config.SetConfigType("yml")
	config.SetConfigName("default")
	config.AddConfigPath("../../configs/")
	config.AddConfigPath(".")

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{config}, nil
}

func (c Config) Sub(key string) ports.Config {
	sub := c.Viper.Sub(key)
	if sub == nil {
		return &Config{viper.New()}
	}
	return &Config{sub}
}

func NewTestConfig(configs ...map[string]any) ports.Config {
	v := viper.New()
	for _, c := range configs {
		err := v.MergeConfigMap(c)
		if err != nil {
			panic(err)
		}
	}
	return &Config{v}
}
