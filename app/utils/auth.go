package utils

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

const (
	HeaderBearerType  = "Bearer"
	HeaderBasicType   = "Basic"
)

func CompareHashAndPasswordFromAuthHeader(passwordHash []byte, r *http.Request) error {
	authToken := r.Header.Get("Authorization")

	encodedPassword, err := GetTokenFromAuthHeader(authToken, HeaderBasicType)
	if err != nil {
		return err
	}

	err = ComparePasswordAndHash(passwordHash, encodedPassword)
	if err != nil {
		return err
	}

	return nil
}

func ComparePasswordAndHash(passwordHash []byte, encodedPassword string) error {
	password, err := base64.StdEncoding.DecodeString(encodedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, password)
	if err != nil {
		return err
	}

	return nil
}

func GetTokenFromAuthHeader(authHeader string, headerType string) (string, error) {
	splitHeader := strings.Split(authHeader, headerType + " ")
	if len(splitHeader) != 2 {
		return "", errors.New("header not valid")
	}
	return splitHeader[1], nil
}
