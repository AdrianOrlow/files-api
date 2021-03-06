package model

import "github.com/AdrianOrlow/files-api/app/utils"

type File struct {
	Model
	Title            string `json:"title"`
	Description      string `json:"description"`
	Password         string `json:"-"`
	HasPassword      bool   `json:"hasPassword" sql:"-"`
	Permalink        string `json:"permalink"`
	FileName         string `json:"fileName"`
	FileSizeKB       string `json:"fileSizeKB"`
	FileChecksumMD5  string `json:"fileChecksumMd5"`
	FileChecksumSHA1 string `json:"fileChecksumSha1"`
	FolderID         int    `json:"-"`
	FolderHashID     string `json:"folderId" sql:"-"`
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
	if f.FolderHashID == "public" {
		f.FolderID = 1
		return f
	}

	f.FolderID, _ = utils.DecodeId(f.FolderHashID, utils.FoldersResourceType)
	return f
}

func (f *File) WithHasPassword() *File {
	if f.Password != "" {
		f.HasPassword = true
	}
	return f
}

func (f *File) WithFileNameWithoutTimestamp() *File {
	f.FileName = utils.GetFileNameWithoutTimestamp(f.FileName)
	return f
}
