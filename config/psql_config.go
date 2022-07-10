package config

import (
	"gin-chat-svc/pkg/logger"
	"os"

	"github.com/spf13/viper"
)

type PsqlConfig struct {
	DBHost			string		`mapstructure:"DBHOST"`
	DBName			string		`mapstructure:"DBNAME"`
	DBUser			string		`mapstructure:"DBUSER"`
	DBPassword		string		`mapstructure:"DBPASSWORD"`
	DBPort			string		`mapstructure:"DBPORT"`
}

func GetPsqlConfig() (conf PsqlConfig, err error) {
	mode := os.Getenv("")

	viper.AddConfigPath(".")
	if mode == "prod" {
		viper.SetConfigFile(".env")
	} else {
		viper.SetConfigName("config")
	}
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		logger.Logger.Error("config", logger.String("error", err.Error()))
	}

	viper.Unmarshal(&conf)

	return
}