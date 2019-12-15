package app

import (
	"fmt"
	"github.com/AdrianOrlow/files-api/app/model"
	"github.com/AdrianOrlow/files-api/config"
	"github.com/jinzhu/gorm"
	"time"
)

func (a *App) InitializeDatabase(config *config.Config) error {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		return err
	}
	a.DB = model.DBMigrate(db)
	a.DB = createPublicFolder(a.DB)

	return nil
}

func createPublicFolder(db *gorm.DB) *gorm.DB {
	var folder model.Folder
	publicFolder := model.Folder{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Title:     "Public",
		Permalink: "public",
		IsPublic:  true,
		ParentID:  0,
	}

	publicFolderRecord := db.First(&folder, 1)
	if publicFolderRecord.RecordNotFound() {
		db.FirstOrCreate(&folder, &publicFolder)
	}

	return db
}
