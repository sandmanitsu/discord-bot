package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	Application Appliction `yaml:"application" env-required:"true"`
	Model       Model      `yaml:"model" env-required:"true"`
}

type Appliction struct {
	BotToken string `yaml:"api_key" env-required:"true"`
	AppID    string `yaml:"app_id" env-required:"true"`
	GuildID  string `yaml:"guild_id" env-required:"true"`
}

type Model struct {
	Token    string `yaml:"api_key" env-required:"true"`
	ModelURL string `yaml:"api_url" env-required:"true"`
}

func MustLoad() *Config {
	configPath := "./config/local.yaml"

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}

	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	return &config
}
