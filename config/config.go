package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_USERNAME  string `mapstructure:"DB_USERNAME"`
	DB_PASSWORD  string `mapstructure:"DB_PASSWORD"`
	DB_HOST      string `mapstructure:"DB_HOST"`
	DB_NAME      string `mapstructure:"DB_NAME"`
	MIGRATE_PATH string `mapstructure:"MIGRATE_PATH"`

	SERVER_HOST string `mapstructure:"SERVER_HOST"`

	SALT string `mapstructure:"SALT"`
}

func New() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigType("app")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return config, nil
}

func (c *Config) GetDBUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", c.DB_USERNAME, c.DB_PASSWORD, c.DB_HOST, c.DB_NAME)
}
