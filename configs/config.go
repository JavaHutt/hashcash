package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config is a configuration object
type Config struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	HashBits          int           `mapstructure:"hash_bits"`
	HashMaxIterations int           `mapstructure:"hash_max_iterations"`
	HashCounter       int           `mapstructure:"hash_counter"`
	HashExpiration    time.Duration `mapstructure:"hash_expiration"`
}

// ParseConfig parses the `config.yaml` file in path
// if path is empty, default config is loaded
func ParseConfig(path string) (*Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if path != "" {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
