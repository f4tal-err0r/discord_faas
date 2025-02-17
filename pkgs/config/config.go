package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Discord   Discord `mapstructure:"discord"`
	Filestore string  `mapstructure:"filestore"`
	DBPath    string  `mapstructure:"dbpath"`
}

type Discord struct {
	Token    string `mapstructure:"token"`
	ClientID string `mapstructure:"clientid"`
	AdminUID string `mapstructure:"adminuid"`
}

func New() (*Config, error) {
	return NewPathConfig("/app/config/config.yaml")
}

func NewPathConfig(path string) (*Config, error) {
	var cfg Config

	if path == "" {
		path = "config.yaml"
	}

	viper.SetConfigFile(path)

	viper.AutomaticEnv()
	viper.BindEnv("discord.token", "DISCORD_TOKEN")

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

	if cfg.Discord.Token == "" {
		return nil, fmt.Errorf("missing config: discord.token")
	}
	if cfg.Discord.ClientID == "" {
		return nil, fmt.Errorf("missing config: discord.clientid")
	}
	return &cfg, nil
}
