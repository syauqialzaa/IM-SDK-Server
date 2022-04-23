package config

import (
	"gin-chat-svc/pkg/logger"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppName			string		`mapstructure:"APP_NAME"`
	ElephantSQL		string		`mapstructure:"ELEPHANT_SQL"`
	LogPath			string		`mapstructure:"LOG_PATH"`
	LogLevel		string		`mapstructure:"LOG_LEVEL"`
	StaticFile		string		`mapstructure:"STATIC_FILE"`
	ChannelType		string		`mapstructure:"CHANNEL_TYPE"`
	KafkaHost		string		`mapstructure:"KAFKA_HOST"`
	KafkaTopic		string		`mapstructure:"KAFKA_TOPIC"`
}

func GetConfig() (conf Config, err error) {
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


