package adapters

import (
	"fmt"
	"os"
	"strings"
	"xmpp-llm-bridge/internal/ports"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

var _ ports.Config = &Config{}

// New returns a new config instance.
func NewConfig() (*Config, error) {
	if err := loadEnvFiles(); err != nil {
		return nil, fmt.Errorf("error loading env files: %w", err)
	}
	config := viper.New()
	config.SetConfigType("yml")
	config.SetConfigName("default")
	config.AddConfigPath("./configs")
	config.AddConfigPath(".")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{config}, nil
}

func loadEnvFiles() error {
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		if !strings.HasSuffix(key, "_FILE") {
			continue
		}

		target := strings.TrimSuffix(key, "_FILE")
		filename := strings.TrimSpace(val)
		if filename == "" {
			continue
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error reading env file %s: %w", filename, err)
		}

		value := strings.TrimRight(string(data), "\r\n")

		if err := os.Setenv(target, value); err != nil {
			return fmt.Errorf("error setting env %s: %w", target, err)
		}
		if err := os.Unsetenv(key); err != nil {
			return fmt.Errorf("error unsetting env %s: %w", key, err)
		}
	}
	return nil
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
