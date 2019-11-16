package utils

import (
	"github.com/AdrianOrlow/files-api/config"
	"github.com/speps/go-hashids"
)

type Utils struct {
	hashID *hashids.HashID
	jwt    JWT
	files  Files
}

type JWT struct {
	secretKey    []byte
	adminsEmails []string
}

type Files struct {
	dir string
}

var utils Utils

func Initialize(config *config.Config) error {
	InitializeJWT(config)

	err := InitializeFiles(config)
	if err != nil {
		return err
	}

	err = InitializeHashId(config)
	return err
}
