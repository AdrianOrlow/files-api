package utils

import (
	"errors"
	"github.com/AdrianOrlow/files-api/config"
	"github.com/dgrijalva/jwt-go"
	"net/http"
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
	utils.jwt.adminsEmails = config.AdminsGMailAddresses
}

func CreateLoginJWT(email string) (*TokenResponse, error) {
	if !loginEmailValid(email) {
		return &TokenResponse{}, errors.New("email not acceptable")
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

	tokenResponse := &TokenResponse{Token: HeaderBearerType + " " + tokenString}

	return tokenResponse, err
}

func VerifyLoginJWT(authHeader string) (int, error) {
	token, err := GetTokenFromAuthHeader(authHeader, HeaderBearerType)
	if err == jwt.ErrSignatureInvalid {
		return http.StatusBadRequest, err
	}

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
		if loginEmailValid(claims["email"]) {
			return http.StatusOK, nil
		}
	}

	return http.StatusUnauthorized, err
}

func loginEmailValid(email interface{}) bool {
	for _, a := range utils.jwt.adminsEmails {
		if a == email {
			return true
		}
	}
	return false
}
