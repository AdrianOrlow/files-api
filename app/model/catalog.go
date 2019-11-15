package model

import "github.com/AdrianOrlow/files-api/app/utils"

type Catalog struct {
	Model
	Title        string `json:"title"`
	Permalink    string `json:"permalink"`
	IsPublic     bool   `json:"is_public"`
	ParentID     int    `json:"-"`
	ParentHashID string `sql:"-" json:"parent_id"`
}

type CatalogPathElement struct {
	Index    int    `json:"index"`
	HashId   string `sql:"-" json:"id"`
	Title    string `json:"title"`
	IsPublic bool   `json:"is_public"`
}

func (c *Catalog) WithHashId() *Catalog {
	c.HashId = utils.EncodeId(int(c.ID), utils.CatalogsResourceType)
	c.ParentHashID = utils.EncodeId(c.ParentID, utils.CatalogsResourceType)
	return c
}

func (c *Catalog) WithId() *Catalog {
	c.ParentID, _ = utils.DecodeId(c.ParentHashID, utils.CatalogsResourceType)
	return c
}

func (c *Catalog) ToPath(index int) CatalogPathElement {
	return CatalogPathElement{
		Index:    index,
		HashId:   c.HashId,
		Title:    c.Title,
		IsPublic: c.IsPublic,
	}
}
