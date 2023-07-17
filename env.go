package chatify

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	SpotifyIDKeyName     = "spotify_id"
	SpotifySecretKeyName = "spotify_secret"
	OpenAIApiKeyName     = "openai_api_key"
	AccessTokenKeyName   = "access_token"
	RefreshTokenKeyName  = "refresh_token"
	ExpirationKeyName    = "expiration"
)

var (
	tokenViper        *viper.Viper
	clientViper       *viper.Viper
	ErrConfigNotFound = errors.New("config file does not exist")
)

func initEnv() error {
	// set config path
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("fail to get user home dir: %w", err)
	}
	configPath := filepath.Join(userHomePath, ".config", "chatify")

	// path check
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return fmt.Errorf("fail to make config path: %v", err)
		}
	}

	// initialize vipers
	tokenViper = viper.New()
	clientViper = viper.New()

	// token setting
	tokenViper.AddConfigPath(configPath)
	tokenViper.SetConfigName("token")
	tokenViper.SetConfigType("json")
	if err := tokenViper.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return fmt.Errorf("fail to safe write token config: %v", err)
		}
	}

	// client setting
	clientViper.AddConfigPath(configPath)
	clientViper.SetConfigName("client")
	clientViper.SetConfigType("yaml")
	if err := clientViper.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return fmt.Errorf("fail to safe write client config: %v", err)
		}
	}

	// load
	if err := tokenViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrConfigNotFound
		} else {
			return fmt.Errorf("fail to read config: %v", err)
		}
	}
	if err := clientViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrConfigNotFound
		} else {
			return fmt.Errorf("fail to read config: %v", err)
		}
	}

	return nil
}

func setTokenInfo(token *oauth2.Token) error {
	tokenViper.Set(AccessTokenKeyName, token.AccessToken)
	tokenViper.Set(RefreshTokenKeyName, token.RefreshToken)
	tokenViper.Set(ExpirationKeyName, token.Expiry.Unix())

	if err := tokenViper.WriteConfig(); err != nil {
		return fmt.Errorf("fail to write token config: %v", err)
	}

	return nil
}

func getToken() *oauth2.Token {
	accessToken := tokenViper.GetString(AccessTokenKeyName)
	refreshToken := tokenViper.GetString(RefreshTokenKeyName)
	expiration := tokenViper.GetString(ExpirationKeyName)
	unix, _ := strconv.ParseInt(expiration, 10, 64)
	expiry := time.Unix(unix, 0)

	return &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}
}

func askClientInfo() error {
	// get spotify ID
	fmt.Println("Chatify requires <Spotify ID>, <Spotify Secret> and <OpenAI API key>")
	fmt.Printf("\nPlease enter your <Spotify ID>\n")
	fmt.Printf("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	getSpotifyID := scanner.Text()

	// get spotify secret
	fmt.Println("Please enter your <Spotify Secret>")
	fmt.Printf("> ")
	scanner.Scan()
	getSpotifySecret := scanner.Text()

	// get OpenAI API key
	fmt.Println("Please enter your <OpenAI API key>")
	fmt.Printf("> ")
	scanner.Scan()
	getOpenAIApiKey := scanner.Text()

	if err := setClientViper(getSpotifyID, getSpotifySecret, getOpenAIApiKey); err != nil {
		return err
	}

	return nil
}

func setClientViper(spotifyID, spotifySecret, openAIApiKey string) error {
	clientViper.Set(SpotifyIDKeyName, spotifyID)
	clientViper.Set(SpotifySecretKeyName, spotifySecret)
	clientViper.Set(OpenAIApiKeyName, openAIApiKey)

	if err := clientViper.WriteConfig(); err != nil {
		return fmt.Errorf("fail to write client config: %v", err)
	}

	return nil
}

func isClientInfoSet() bool {
	spotifyID := clientViper.GetString(SpotifyIDKeyName)
	spotifySecret := clientViper.GetString(SpotifySecretKeyName)
	openAIApiKey := clientViper.GetString(OpenAIApiKeyName)

	return spotifyID != "" && spotifySecret != "" && openAIApiKey != ""
}

func GetOpenAIAPIKey() string {
	return clientViper.GetString(OpenAIApiKeyName)
}
