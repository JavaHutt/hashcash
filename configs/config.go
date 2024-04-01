package configs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config is a configuration object
type Config struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	StoreAddr         string        `mapstructure:"store_addr"`
	StoreExpiration   time.Duration `mapstructure:"store_expiration"`
	HashBits          int           `mapstructure:"hash_bits"`
	HashMaxIterations int           `mapstructure:"hash_max_iterations"`
	HashCounter       int           `mapstructure:"hash_counter"`
	HashExpiration    time.Duration `mapstructure:"hash_expiration"`
}

// ParseConfig parses the `config.yaml` file in path
// if path is empty, default config is loaded
func ParseConfig(path string) (*Config, error) {
	config, err := loadDefaultConfig(path)
	if err != nil {
		return nil, err
	}

	host := viper.GetString("SERVER_HOST")
	if host != "" {
		config.Host = host
	}

	portStr := viper.GetString("SERVER_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		config.Port = port
	}

	storeAddr := viper.GetString("STORE_ADDR")
	if storeAddr != "" {
		config.StoreAddr = storeAddr
	}

	return config, nil
}

func loadDefaultConfig(path string) (*Config, error) {
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

	viper.AutomaticEnv()

	return &config, nil
}
