package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Discord     Discord `mapstructure:"DISCORD"`
	Filestore   string  `mapstructure:"FILESTORE"`
	DBPath      string  `mapstructure:"DBPATH"`
	RuntimeRepo string  `mapstructure:"RUNTIMEREPO"`
}

type Discord struct {
	Token    string `mapstructure:"TOKEN"`
	ClientID string `mapstructure:"CLIENTID"`
}

func New() (*Config, error) {
	return NewPathConfig("")
}

func NewPathConfig(path string) (*Config, error) {
	var cfg Config

	if path == "" {
		path = "config.yaml"
	}

	viper.SetConfigFile(path)

	viper.SetDefault("filestore", "/opt/dfaas")
	viper.SetDefault("runtimerepo", "github.com/f4tal-err0r/discord_faas/runtimes/")
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
