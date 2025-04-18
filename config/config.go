package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDB MongoDBConfig
	Kafka   KafkaConfig
	Server  ServerConfig
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into config struct: %v", err)
	}

	return &config
}