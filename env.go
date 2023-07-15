package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	configPath = "$HOME/.config/chatify"
)

func setEnv() error {
	// token setting
	viper.AddConfigPath(configPath)
	viper.SetConfigName("token")
	viper.SetConfigType("json")

	// client setting
	viper.SetConfigName("client")
	viper.SetConfigType("yaml")
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("fail to merge config: %v", err)
	}

	// load
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("config file does not exist")
		}
	}

	return nil
}
