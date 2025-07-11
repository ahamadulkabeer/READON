package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBUrl                string `mapstructure:"DB_URL"`
	RazorpayKey          string `mapstructure:"RAZORPAY_SECRET"`
	RazorpaySecret       string `mapstructure:"RAZORPAY_KEY"`
	EmailjetApiKey       string `mapstructure:"EMAILJET_KEY"`
	EmailjetSecretKey    string `mapstructure:"EMAILJET_SECRET"`
	AWSS3AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSS3SecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	JWTSecretKeyword     string `mapstructure:"JWT_SECRET_KEYWORD"`
}

var envs = []string{
	"DB_URL",
	"RAZORPAY_SECRET", "RAZORPAY_KEY", "EMAILJET_KEY", "EMAILJET_SECRET",
	"AWS_SECRET_ACCESS_KEY", "AWS_ACCESS_KEY_ID", "JWT_SECRET_KEYWORD",
}

func LoadConfig() (Config, error) {
	var config Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./.env" // to run in localhost
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

	return config, nil
}
