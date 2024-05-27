package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Discord Discord `mapstructure:"DISCORD"`
	Domain  string  `mapstructure:"DOMAIN"`
	Oauth   Oauth   `mapstructure:"OAUTH"`
}

type Discord struct {
	Token   string   `mapstructure:"BOT_TOKEN"`
	AdminId []string `mapstructure:"ADMIN_ID"`
}

type Oauth struct {
	ClientID     string `mapstructure:"CLIENTID"`
	ClientSecret string `mapstructure:"CLIENTSECRET"`
}

func New() (*Config, error) {
	var config Config

	// Set up Viper to read the config.yaml file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Check if the error is due to the file not existing
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %s", err)
		}
	}
	// Set up Viper to read secret environment variables
	viper.BindEnv("oauth.clientid", "DFAAS_OAUTH_CLIENTID")
	viper.BindEnv("oauth.clientsecret", "DFAAS_OAUTH_CLIENTSECRET")
	viper.BindEnv("discord.token", "DFAAS_OAUTH_CLIENTID")
	viper.BindEnv("oauth.clientsecret", "DFAAS_OAUTH_CLIENTSECRET")

	// Unmarshal the configuration into the struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *Config) FetchCache() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Unable to fetch cache directory: %v", err)
	}
	cacheDir := cache + "/dfaas"
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		log.Fatalf("Unable to create cache directory: %v", err)
	}
	return filepath.Join(cacheDir, "token.json")
}
