package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost            string `mapstructure:"DB_HOST"`
	DBName            string `mapstructure:"DB_NAME"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	RazorpayKey       string `mapstructure:"RAZORPAY_SECRET"`
	RazorpaySecret    string `mapstructure:"RAZORPAY_KEY"`
	EmailjetApiKey    string `mapstructure:"EMAILJET_KEY"`
	EmailjetSecretKey string `mapstructure:"EMAILJET_SECRET"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "KEY", "RAZORPAY_SECRET", "RAZORPAY_KEY", "EMAILJET_KEY", "EMAILJET_SECRET",
}

func LoadConfig() (Config, error) {
	var config Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./.env"
	}

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}
	fmt.Println("config  : ", config)

	return config, nil
}
