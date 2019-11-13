package utils

import (
	"errors"
	"github.com/AdrianOrlow/files-api/config"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Token string `json:"token"`
}

func InitializeJWT(config *config.Config) {
	utils.jwt.secretKey = []byte(config.SecretJWT)
	utils.jwt.adminsEmails = config.AdminsGmailAddresses
}

func CreateJWT(email string) (*TokenResponse, error) {
	if !emailValid(email) {
		return &TokenResponse{}, errors.New("EmailNotAcceptable")
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.jwt.secretKey)

	tokenResponse := &TokenResponse{Token: "Bearer " + tokenString}

	return tokenResponse, err
}

func VerifyJWT(token string) (int, error) {
	token = strings.Split(token, "Bearer ")[1]

	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return utils.jwt.secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, err
		}
		return http.StatusBadRequest, err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		if emailValid(claims["email"]) {
			return http.StatusOK, nil
		}
	}

	return http.StatusUnauthorized, err
}

func emailValid(email interface{}) bool {
	for _, a := range utils.jwt.adminsEmails {
		if a == email {
			return true
		}
	}
	return false
}