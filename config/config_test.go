package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestConfKey_isTokenKey(t *testing.T) {
	type fields struct {
		confkey ConfKey
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "true; access token",
			fields: fields{
				confkey: AccessTokenKey,
			},
			want: true,
		},
		{
			name: "true; refresh token",
			fields: fields{
				confkey: RefreshTokenKey,
			},
			want: true,
		},
		{
			name: "true; expiration",
			fields: fields{
				confkey: ExpirationKey,
			},
			want: true,
		},
		{
			name: "false; spotify id",
			fields: fields{
				confkey: SpotifyIDKey,
			},
			want: false,
		},
		{
			name: "false; spotify secret",
			fields: fields{
				confkey: SpotifySecretKey,
			},
			want: false,
		},
		{
			name: "false; openai api",
			fields: fields{
				confkey: OpenAIAPIKey,
			},
			want: false,
		},
		{
			name: "false; device id",
			fields: fields{
				confkey: DeviceID,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.confkey.isTokenKey() != tt.want {
				t.Errorf("got: %v, want: %v", tt.fields.confkey.isTokenKey(), tt.want)
			}
		})
	}
}

func TestConfKey_isClientKey(t *testing.T) {
	type fields struct {
		confkey ConfKey
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "false; access token",
			fields: fields{
				confkey: AccessTokenKey,
			},
			want: false,
		},
		{
			name: "false; refresh token",
			fields: fields{
				confkey: RefreshTokenKey,
			},
			want: false,
		},
		{
			name: "false; expiration",
			fields: fields{
				confkey: ExpirationKey,
			},
			want: false,
		},
		{
			name: "true; spotify id",
			fields: fields{
				confkey: SpotifyIDKey,
			},
			want: true,
		},
		{
			name: "true; spotify secret",
			fields: fields{
				confkey: SpotifySecretKey,
			},
			want: true,
		},
		{
			name: "true; openai api",
			fields: fields{
				confkey: OpenAIAPIKey,
			},
			want: true,
		},
		{
			name: "true; device id",
			fields: fields{
				confkey: DeviceID,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.confkey.isClientKey() != tt.want {
				t.Errorf("got: %v, want: %v", tt.fields.confkey.isClientKey(), tt.want)
			}
		})
	}
}

func Test_constructConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigDir)
	if err := constructConfigPath(configPath); err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("fail to construct config path")
	}
}

func TestSetupViper(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigDir)
	if err := constructConfigPath(configPath); err != nil {
		t.Errorf("fail to construct config path: %v", err)
		return
	}
	if err := setupViper(viper.New(), configPath, "token", "json"); err != nil {
		t.Errorf("fail to setup token viper: %v", err)
	}
	if err := setupViper(viper.New(), configPath, "client", "yaml"); err != nil {
		t.Errorf("fail to setup client viper: %v", err)
	}
}

const testValue = "test value"

func TestConfig_Load(t *testing.T) {
	conf := newConfig()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigDir)
	if err := constructConfigPath(configPath); err != nil {
		t.Errorf("fail to construct config path: %v", err)
		return
	}
	if err := setupViper(conf.tokenViper, configPath, "token", "json"); err != nil {
		t.Errorf("fail to setup token viper: %v", err)
	}
	if err := setupViper(conf.clientViper, configPath, "client", "yaml"); err != nil {
		t.Errorf("fail to setup client viper: %v", err)
	}

	f, err := os.Create(filepath.Join(configPath, "token.json"))
	if err != nil {
		t.Errorf("file open error: %v", err)
	}
	if _, err := io.WriteString(f, fmt.Sprintf(`{"test": "%s"}`, testValue)); err != nil {
		t.Errorf("fail to write token.json: %v", err)
	}
	f.Close()

	f, err = os.Create(filepath.Join(configPath, "client.yaml"))
	if err != nil {
		t.Errorf("file open error: %v", err)
	}
	if _, err := io.WriteString(f, fmt.Sprintf(`test: %s`, testValue)); err != nil {
		t.Errorf("fail to write client.yaml: %v", err)
	}
	f.Close()

	if err := conf.Load(); err != nil {
		t.Errorf("fail to load: %v", err)
	}

	if conf.tokenViper.GetString("test") != testValue {
		t.Errorf("got: %s, want: %s", conf.tokenViper.GetString("test"), testValue)
	}
	if conf.clientViper.GetString("test") != testValue {
		t.Errorf("got: %s, want: %s", conf.clientViper.GetString("test"), testValue)
	}
}

func TestConf_IsClientValid(t *testing.T) {
	tests := []struct {
		name          string
		spotifyID     string
		spotifySecret string
		openAIapiKey  string
		deviceID      string
		want          bool
	}{
		{
			name:          "true",
			spotifyID:     "test",
			spotifySecret: "test",
			openAIapiKey:  "test",
			deviceID:      "test",
			want:          true,
		},
		{
			name:          "false",
			spotifyID:     "",
			spotifySecret: "",
			openAIapiKey:  "",
			deviceID:      "",
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := newConfig()

			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ConfigDir)
			if err := constructConfigPath(configPath); err != nil {
				t.Errorf("fail to construct config path: %v", err)
				return
			}
			if err := setupViper(conf.tokenViper, configPath, "token", "json"); err != nil {
				t.Errorf("fail to setup token viper: %v", err)
			}
			if err := setupViper(conf.clientViper, configPath, "client", "yaml"); err != nil {
				t.Errorf("fail to setup client viper: %v", err)
			}

			conf.clientViper.Set(string(SpotifyIDKey), tt.spotifyID)
			conf.clientViper.Set(string(SpotifySecretKey), tt.spotifySecret)
			conf.clientViper.Set(string(OpenAIAPIKey), tt.openAIapiKey)
			conf.clientViper.Set(string(DeviceID), tt.deviceID)

			if conf.IsClientValid() != tt.want {
				t.Errorf("got: %v, want: %v", conf.IsClientValid(), tt.want)
			}
		})
	}
}

func TestConf_GetClientValue(t *testing.T) {
	type args struct {
		key ConfKey
	}
	tests := []struct {
		name string
		args
		clientKey ConfKey
		want      string
	}{
		{
			name: "success: client key",
			args: args{
				key: SpotifyIDKey,
			},
			clientKey: SpotifyIDKey,
			want:      testValue,
		},
		{
			name: "success: token key",
			args: args{
				key: AccessTokenKey,
			},
			clientKey: SpotifyIDKey,
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := newConfig()

			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ConfigDir)
			if err := constructConfigPath(configPath); err != nil {
				t.Errorf("fail to construct config path: %v", err)
				return
			}
			if err := setupViper(conf.tokenViper, configPath, "token", "json"); err != nil {
				t.Errorf("fail to setup token viper: %v", err)
			}
			if err := setupViper(conf.clientViper, configPath, "client", "yaml"); err != nil {
				t.Errorf("fail to setup client viper: %v", err)
			}

			conf.clientViper.Set(string(tt.clientKey), testValue)

			if err := conf.Load(); err != nil {
				t.Errorf("fail to load: %v", err)
			}

			if conf.GetClientValue(tt.args.key) != tt.want {
				t.Errorf("got: %s, want: %s", conf.GetClientValue(tt.args.key), tt.want)
			}
		})
	}
}
