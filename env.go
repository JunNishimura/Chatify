package main

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	configPath = "$HOME/.config/chatify"
)

func setEnv() {
	// token setting
	viper.AddConfigPath(configPath)
	viper.SetConfigName("token")
	viper.SetConfigType("json")

	// client setting
	viper.SetConfigName("client")
	viper.SetConfigType("yaml")
	viper.MergeInConfig()

	// load
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config file does not exist")
		}
	}
}
