package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Discord *Discord
	Domain  string
	Docker  *Docker
}

type Discord struct {
	Token     string
	AdminId   string   //Discord ID of admin
	Approvers []string //DiscordID of valid approvers
}

type Docker struct {
	Registry string
	Limit    int //Limit to number of images per Discord server
}

func New(path string) *Config {
	var c *Config

	viper.SetEnvPrefix("DFAAS")
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("ERROR: Unable to read in config")
	}

	// Unmarshal configuration into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("ERROR: Unable to marshal config")
	}

	return c
}
