package model

import "github.com/AdrianOrlow/files-api/app/utils"

type Folder struct {
	Model
	Title        string `json:"title"`
	Permalink    string `json:"permalink"`
	IsPublic     bool   `json:"is_public"`
	ParentID     int    `json:"-"`
	ParentHashID string `sql:"-" json:"parent_id"`
}

type FolderPathElement struct {
	Index    int    `json:"index"`
	HashId   string `sql:"-" json:"id"`
	Title    string `json:"title"`
	IsPublic bool   `json:"is_public"`
}

func (f *Folder) WithHashId() *Folder {
	f.HashId = utils.EncodeId(int(f.ID), utils.FoldersResourceType)
	f.ParentHashID = utils.EncodeId(f.ParentID, utils.FoldersResourceType)
	return f
}

func (f *Folder) WithId() *Folder {
	f.ParentID, _ = utils.DecodeId(f.ParentHashID, utils.FoldersResourceType)
	return f
}

func (c *Folder) ToPath(index int) FolderPathElement {
	return FolderPathElement{
		Index:    index,
		HashId:   c.HashId,
		Title:    c.Title,
		IsPublic: c.IsPublic,
	}
}
