package model

import "github.com/AdrianOrlow/files-api/app/utils"

type Folder struct {
	Model
	Title        string `json:"title"`
	Permalink    string `json:"permalink"`
	ParentID     int    `json:"-"`
	ParentHashID string `sql:"-" json:"parentId"`
}

type FolderPathElement struct {
	Index    int    `json:"index"`
	HashId   string `sql:"-" json:"id"`
	Title    string `json:"title"`
}

func (f *Folder) WithHashId() *Folder {
	f.HashId = utils.EncodeId(int(f.ID), utils.FoldersResourceType)
	f.ParentHashID = utils.EncodeId(f.ParentID, utils.FoldersResourceType)
	return f
}

func (f *Folder) WithId() *Folder {
	if f.ParentHashID == "public" {
		f.ParentID = 1
		return f
	}

	f.ParentID, _ = utils.DecodeId(f.ParentHashID, utils.FoldersResourceType)
	return f
}

func (f *Folder) ToPath(index int) FolderPathElement {
	return FolderPathElement{
		Index:    index,
		HashId:   f.HashId,
		Title:    f.Title,
	}
}
