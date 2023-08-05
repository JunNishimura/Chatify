package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	ConfigDir = ".config/chatify"
)

type Config struct {
	tokenViper  *viper.Viper
	clientViper *viper.Viper
}

type ConfKey string

const (
	SpotifyIDKey     ConfKey = "spotify_id"
	SpotifySecretKey ConfKey = "spotify_secret"
	OpenAIAPIKey     ConfKey = "openai_api_key"
	AccessTokenKey   ConfKey = "access_token"
	RefreshTokenKey  ConfKey = "refresh_token"
	ExpirationKey    ConfKey = "expiration"
)

func Load() (*Config, error) {
	conf := new(Config)

	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if err := constructConfigPath(configPath); err != nil {
		return nil, err
	}

	if err := setupViper(conf.tokenViper, configPath, "token", "json"); err != nil {
		return nil, err
	}
	if err := setupViper(conf.clientViper, configPath, "client", "yaml"); err != nil {
		return nil, err
	}

	if err := conf.tokenViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fail to read token viper: %v", err)
	}
	if err := conf.clientViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fail to read client viper: %v", err)
	}

	return conf, nil
}

func getConfigPath() (string, error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("fail to get user home dir: %v", err)
	}

	return filepath.Join(d, ConfigDir), nil
}

func constructConfigPath(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return fmt.Errorf("fail to make config path(%s): %v", configPath, err)
		}
	}

	return nil
}

func setupViper(v *viper.Viper, configPath, configName, configType string) error {
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	if err := v.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return fmt.Errorf("fail to write config(%s/%s.%s): %v",
				configPath,
				configName,
				configType,
				err,
			)
		}
	}

	return nil
}

func (c *Config) SetToken(token *oauth2.Token) error {
	c.tokenViper.Set(string(AccessTokenKey), token.AccessToken)
	c.tokenViper.Set(string(RefreshTokenKey), token.RefreshToken)
	c.tokenViper.Set(string(ExpirationKey), token.Expiry.Unix())

	if err := c.tokenViper.WriteConfig(); err != nil {
		return fmt.Errorf("fail to write token config: %v", err)
	}

	return nil
}

func (c *Config) GetToken() *oauth2.Token {
	accessToken := c.tokenViper.GetString(string(AccessTokenKey))
	refreshToken := c.tokenViper.GetString(string(RefreshTokenKey))
	expiration := c.tokenViper.GetString(string(ExpirationKey))
	unix, _ := strconv.ParseInt(expiration, 10, 64)
	expiry := time.Unix(unix, 0)

	return &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}
}

func (c *Config) SetClient(spotifyID, spotifySecret, openAIApiKey string) error {
	c.clientViper.Set(string(spotifyID), spotifyID)
	c.clientViper.Set(string(spotifySecret), spotifySecret)
	c.clientViper.Set(string(openAIApiKey), openAIApiKey)

	if err := c.clientViper.WriteConfig(); err != nil {
		return fmt.Errorf("fail to write client config: %v", err)
	}

	return nil
}

func (c *Config) IsClientValid() bool {
	spotifyID := c.clientViper.GetString(string(SpotifyIDKey))
	spotifySecret := c.clientViper.GetString(string(SpotifySecretKey))
	openAIApiKey := c.clientViper.GetString(string(OpenAIAPIKey))

	return spotifyID != "" && spotifySecret != "" && openAIApiKey != ""
}
