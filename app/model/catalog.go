package model

import "github.com/AdrianOrlow/files-api/app/utils"

type Catalog struct {
	Model
	Title string `json:"title"`
	Permalink string `json:"permalink"`
	IsPublic bool `json:"is_public"`
	ParentID int `json:"-"`
	ParentHashID string `json:"parent_id" sql:"-"`
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
