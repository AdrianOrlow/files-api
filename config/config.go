package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DB                   *DBConfig
	HashID               *HashID
	GoogleOauthConfig    *GoogleOauthConfig
	AdminsGMailAddresses []string
	StorageDir           string
	SecretJWT            string
	Port                 string
}

type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Host     string
	Port     int
	Name     string
	Charset  string
}

type HashID struct {
	Salt      string
	MinLength int
}

type GoogleOauthConfig struct {
	ClientID     string
	ClientSecret string
}

func LoadConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  getEnv("DB_DIALECT"),
			Username: getEnv("DB_USERNAME"),
			Password: getEnv("DB_PASSWORD"),
			Host:     getEnv("DB_HOST"),
			Port:     getEnvAsInt("DB_PORT"),
			Name:     getEnv("DB_NAME"),
			Charset:  getEnv("DB_CHARSET"),
		},
		HashID: &HashID{
			Salt: getEnv("HASH_ID_SALT"),
			MinLength: getEnvAsInt("HASH_ID_MIN_LENGTH"),
		},
		GoogleOauthConfig: &GoogleOauthConfig{
			ClientID:     getEnv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: getEnv("GOOGLE_OAUTH_CLIENT_SECRET"),
		},
		AdminsGMailAddresses: getEnvAsSlice("ADMIN_GMAIL_ADDRESSES", ","),
		StorageDir: getEnv("STORAGE_DIR"),
		SecretJWT: getEnv("SECRET_JWT"),
		Port: getEnv("PORT"),
	}
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		log.Fatal("ENV VARIABLE DOESN'T EXISTS: " + key)
	}

	return value
}

func getEnvAsInt(key string) int {
	valueStr := getEnv(key)
	value, err := strconv.Atoi(valueStr)

	if err != nil {
		log.Fatal("ENV VARIABLE ISN'T INTEGER: " + key)
	}

	return value
}

func getEnvAsSlice(key string, separator string) []string {
	valStr := getEnv(key)
	val := strings.Split(valStr, separator)

	return val
}
