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
	Storage   Storage `mapstructure:"storage"`
}

type Discord struct {
	Token    string `mapstructure:"token"`
	ClientID string `mapstructure:"clientid"`
}

type Storage struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
	S3   struct {
		Hostname string `mapstructure:"hostname"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"s3"`
}

func New(path string) (*Config, error) {
	return NewPathConfig(path)
}

func NewPathConfig(path string) (*Config, error) {
	var cfg Config

	if path == "" {
		path = "config.yaml"
	}

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	cpath, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("unable to resolve CacheDir: %s", err)
	}
	viper.SetDefault("cachepath", cpath)
	viper.SetDefault("filestore", "/app/data/artifacts")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.Discord.Token = viper.GetString("DISCORD_TOKEN")
	cfg.Discord.ClientID = viper.GetString("DISCORD_CLIENTID")

	requiredConfigs := map[string]*string{
		"discord.token":    &cfg.Discord.Token,
		"discord.clientid": &cfg.Discord.ClientID,
		"storage.path":     &cfg.Storage.Path,
	}
	if cfg.Storage.Type == "s3" {
		requiredConfigs["storage.s3.hostname"] = &cfg.Storage.S3.Hostname
		requiredConfigs["storage.s3.username"] = &cfg.Storage.S3.Username
		requiredConfigs["storage.s3.password"] = &cfg.Storage.S3.Password
	}

	for key, value := range requiredConfigs {
		if *value == "" {
			return nil, fmt.Errorf("missing config: %s", key)
		}
	}

	return &cfg, nil
}
