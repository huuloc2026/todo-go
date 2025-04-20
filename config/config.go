package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Name      string `mapstructure:"name"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parseTime"`
	Loc       string `mapstructure:"loc"`
}

type Config struct {
	DB DBConfig `mapstructure:"db"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// Log ra cấu hình đã được load
	fmt.Printf("Loaded DB Config: %+v\n", AppConfig.DB)
	log.Println("Load config oke")
}
