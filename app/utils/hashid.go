package utils

import (
	"errors"
	"github.com/AdrianOrlow/files-api/config"
	"github.com/speps/go-hashids"
)

const (
	CatalogsResourceType = 0
	FilesResourceType    = 1
)

func InitializeHashId(config *config.Config) error {
	hd := hashids.NewData()
	hd.Salt = config.HashID.Salt
	hd.MinLength = config.HashID.MinLength
	hid, err := hashids.NewWithData(hd)
	utils.hashID = hid
	return err
}

func EncodeId(id int, resourceType int) string {
	e, _ := utils.hashID.Encode([]int{id, resourceType})
	return e
}

func DecodeId(hashId string, resourceType int) (int, error) {
	d, _ := utils.hashID.DecodeWithError(hashId)
	if len(d) != 2 {
		return 0, errors.New("not valid hashId")
	} else if d[1] != resourceType {
		return 0, errors.New("bad resource type")
	} else {
		return d[0], nil
	}
}
