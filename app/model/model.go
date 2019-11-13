package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Model struct {
	ID        uint `gorm:"primary_key" json:"-"`
	HashId string `sql:"-" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Catalog{}, &File{})
	return db
}
