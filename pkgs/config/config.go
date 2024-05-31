package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Discord   Discord `mapstructure:"DISCORD"`
	Domain    string  `mapstructure:"DOMAIN"`
	Filestore string  `mapstructure:"FILESTORE"`
	DBPath    string  `mapstructure:"DBPATH"`
}

type Discord struct {
	Token    string `mapstructure:"TOKEN"`
	ClientID string `mapstructure:"CLIENTID"`
}

func New() (*Config, error) {
	var config Config

	// Set up Viper to read the config.yaml file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("filestore", "/opt/dfaas")

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Check if the error is due to the file not existing
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %s", err)
		}
	}
	// Set up Viper to read secret environment variables
	viper.BindEnv("discord.clientid", "DFAAS_OAUTH_CLIENTID")
	viper.BindEnv("discord.token", "DFAAS_TOKEN")

	// Unmarshal the configuration into the struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}
