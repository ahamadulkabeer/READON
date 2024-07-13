package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBName         string `mapstructure:"DB_NAME"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	RazorpayKey    string `mapstructure:"RAZORPAY_SECRET"`
	RazorpaySecret string `mapstructure:"RAZORPAY_KEY"`
	SendgridApiKey string `mapstructure:"SENDGRID_API_KEY"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "KEY", "RAZORPAY_SECRET", "RAZORPAY_KEY", "SENDGRID_API_KEY",
}

func LoadConfig() (Config, error) {
	var config Config

	// if runningAsSystemdService() {
	// 	viper.SetConfigFile("/home/kabeer/Projects/READON/.env") // Adjust path as needed
	// 	fmt.Println("sytem running >>>>>>>>>>>>")
	// } else {
	// 	viper.AddConfigPath("./")
	// 	viper.SetConfigFile(".env")
	// 	fmt.Println("vc code  running >>>>>>>>>>>>")
	// }
	viper.SetConfigFile("/home/kabeer/Projects/READON/.env")
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

// func runningAsSystemdService() bool {
// 	// Check if the environment variable SYSTEMD_UNIT is set
// 	_, exists := os.LookupEnv("SYSTEMD_UNIT")
// 	return exists
// }
