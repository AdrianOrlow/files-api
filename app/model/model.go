package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	HashId    string     `sql:"-" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Folder{}, &File{})
	return db
}
