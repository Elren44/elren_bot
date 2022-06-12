package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Token         string `mapstructure:"token"`
	VideocdnToken string `mapstructure:"videocdntoken"`
}

func InitConfig() (*Config, error) {
	// viper.AddConfigPath(".")
	// viper.SetConfigType("env")
	// viper.SetConfigFile(".env")
	viper.SetConfigFile("./.env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %w", err)
	}
	return &cfg, nil
}
