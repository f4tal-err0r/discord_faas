package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Discord   Discord `mapstructure:"DISCORD"`
	Domain    string  `mapstructure:"DOMAIN"`
	Filestore string  `mapstructure:"FILESTORE"`
	DBPath    string  `mapstructure:"DBPATH"`
	Cachepath string  `mapstructure:"CACHEPATH"`
}

type Discord struct {
	Token    string `mapstructure:"TOKEN"`
	ClientID string `mapstructure:"CLIENTID"`
}

func New() (*Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("filestore", "/opt/dfaas")
	cpath, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve CacheDir: %s.", err)
	}
	viper.SetDefault("cachepath", cpath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
