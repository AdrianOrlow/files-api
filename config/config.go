package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	DB                   *DBConfig          `json:"db_config"`
	HashID               *HashID            `json:"hash_id"`
	GoogleOauthConfig    *GoogleOauthConfig `json:"google_oauth_config"`
	AdminsGMailAddresses []string           `json:"admins_gmail_addresses"`
	FilesDir             string             `json:"files_dir"`
	SecretJWT            string             `json:"secret_jwt"`
}

type DBConfig struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Charset  string `json:"charset"`
}

type HashID struct {
	Salt      string `json:"salt"`
	MinLength int    `json:"min_length"`
}

type GoogleOauthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func LoadConfig(filepath string) (*Config, error) {
	// Get the config file
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configFile, config)

	return config, err
}
