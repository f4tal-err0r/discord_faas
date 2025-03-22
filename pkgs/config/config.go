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
	URLDomain string  `mapstructure:"domain"`
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

	envBindings := map[string]string{
		"discord.token":    "DISCORD_TOKEN",
		"discord.clientid": "DISCORD_CLIENTID",
		"discord.adminuid": "DISCORD_ADMINUID",
	}
	for key, env := range envBindings {
		viper.BindEnv(key, env)
	}

	cpath, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("unable to resolve CacheDir: %s", err)
	}
	viper.SetDefault("cachepath", cpath)
	viper.SetDefault("filestore", "/app/data/funcstore")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	requiredConfigs := map[string]string{
		"discord.token":    cfg.Discord.Token,
		"discord.clientid": cfg.Discord.ClientID,
	}
	for key, value := range requiredConfigs {
		if value == "" {
			return nil, fmt.Errorf("missing config: %s", key)
		}
	}

	return &cfg, nil
}
