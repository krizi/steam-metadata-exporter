package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type FeatureConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Schedule string `mapstructure:"schedule"`
	APIURL   string `mapstructure:"api_url"`
}

type SteamConfig struct {
	APIKey      string                   `mapstructure:"api_key"`
	UserID      string                   `mapstructure:"user_id"`
	Features    map[string]FeatureConfig `mapstructure:"features"`
	LogLevel    string                   `mapstructure:"log_level"`
	MetricsPort string                   `mapstructure:"metrics_port"`
}

func LoadConfig(logger *logrus.Entry) *SteamConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("STEAM_EXPORTER")

	viper.SetDefault("log_level", "info")
	viper.SetDefault("metrics_port", ":8080")

	if err := viper.ReadInConfig(); err != nil {
		logger.Infof("No configuration file found: %v. Using defaults or environment variables.", err)
	}

	var config SteamConfig
	if err := viper.Unmarshal(&config); err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	return &config
}
