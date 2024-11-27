package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	DSN      string
	DbName   string
	BasePath string
}

func Load(envFile string) *Config {
	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal().Msg("config file not found")
		} else {
			log.Fatal().Err(err).Msg("failed to read config file")
		}
	}

	return &Config{
		DSN:      viper.GetString("DSN"),
		DbName:   viper.GetString("DB_NAME"),
		BasePath: viper.GetString("BASE_PATH"),
	}
}
