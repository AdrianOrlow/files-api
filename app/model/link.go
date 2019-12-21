package model

import (
	"github.com/AdrianOrlow/files-api/app/utils"
	"time"
)

type Link struct {
	Model
	Key        string     `json:"key"`
	FileID     int        `json:"-"`
	FileHashID string     `json:"fileId"`
	ValidUntil time.Time  `json:"validUntil"`
}

func (l *Link) WithHashId() *Link {
	l.HashId = utils.EncodeId(int(l.ID), utils.FilesResourceType)
	l.FileHashID = utils.EncodeId(l.FileID, utils.FilesResourceType)
	return l
}

func (l *Link) IsValid() bool {
	return l.ValidUntil.After(time.Now())
}