package utils

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func CompareHashAndPasswordFromAuthHeader(passwordHash []byte, r *http.Request) error {
	password := getPasswordFromAuthorizationHeader(r)
	if password != nil {
		err := bcrypt.CompareHashAndPassword(passwordHash, password)
		if err != nil {
			return err
		}
	} else {
		err := errors.New("authorization header not provided")
		return err
	}
	return nil
}

func getPasswordFromAuthorizationHeader(r *http.Request) []byte {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return nil
	}

	splitToken := strings.Split(reqToken, "Basic ")
	if len(splitToken) != 2 {
		return nil
	}

	reqToken = splitToken[1]
	password, _ := base64.StdEncoding.DecodeString(reqToken)
	return password
}
