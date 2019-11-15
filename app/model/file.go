package model

import "github.com/AdrianOrlow/files-api/app/utils"

type File struct {
	Model
	Title            string `json:"title"`
	Description      string `json:"description"`
	Password         string `json:"-"`
	HasPassword      bool   `json:"has_password" sql:"-"`
	Permalink        string `json:"permalink"`
	FileName         string `json:"filename"`
	FileSizeKB       string `json:"file_size_kb"`
	FileChecksumMD5  string `json:"file_checksum_md5"`
	FileChecksumSHA1 string `json:"file_checksum_sha1"`
	CatalogID        int    `json:"-"`
	CatalogHashID    string `json:"catalog_id" sql:"-"`
}

func (f *File) WithHashId() *File {
	f.HashId = utils.EncodeId(int(f.ID), utils.FilesResourceType)
	f.CatalogHashID = utils.EncodeId(f.CatalogID, utils.CatalogsResourceType)
	return f
}

func (f *File) WithHasPassword() *File {
	if f.Password != "" {
		f.HasPassword = true
	}
	return f
}
