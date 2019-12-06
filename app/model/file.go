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
	FolderID        int    `json:"-"`
	FolderHashID    string `json:"folder_id" sql:"-"`
}

type FilePassword struct {
	Password string `json:"password"`
}

func (f *File) WithHashId() *File {
	f.HashId = utils.EncodeId(int(f.ID), utils.FilesResourceType)
	f.FolderHashID = utils.EncodeId(f.FolderID, utils.FoldersResourceType)
	return f
}

func (f *File) WithId() *File {
	f.FolderID, _ = utils.DecodeId(f.FolderHashID, utils.FoldersResourceType)
	return f
}

func (f *File) WithHasPassword() *File {
	if f.Password != "" {
		f.HasPassword = true
	}
	return f
}
