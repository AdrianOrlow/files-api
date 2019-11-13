package utils

import (
	"github.com/AdrianOrlow/files-api/config"
	"github.com/speps/go-hashids"
)

type Utils struct {
	hashID *hashids.HashID
	jwt    JWT
}

type JWT struct {
	secretKey    []byte
	adminsEmails []string
}

var utils Utils

func Initialize(config *config.Config) error {
	InitializeJWT(config)

	err := InitializeHashId(config)
	return err
}
