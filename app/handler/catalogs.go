package handler

import (
	"encoding/json"
	"github.com/AdrianOrlow/files-api/app/utils"
	"net/http"

	"github.com/AdrianOrlow/files-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetRootCatalogs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	catalogs := getRootCatalogsOr404(db, w, r)
	if catalogs == nil {
		return
	}
	respondJSON(w, http.StatusOK, catalogs)
}

func CreateCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var catalog model.Catalog

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&catalog); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(catalog.WithId()).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, catalog.WithHashId())
}

func UpdateCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.CatalogsResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		catalog := getCatalogOr404(db, uint(id), w, r)
		if catalog == nil {
			return
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&catalog); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		if err := db.Save(catalog.WithId()).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, catalog.WithHashId())
	}
}

func DeleteCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.CatalogsResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		catalog := getCatalogOr404(db, uint(id), w, r)
		if catalog == nil {
			return
		}
		if err := db.Delete(&catalog).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusNoContent, nil)
	}
}

func GetCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.CatalogsResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		catalog := getCatalogOr404(db, uint(id), w, r)
		if catalog == nil {
			return
		}
		respondJSON(w, http.StatusOK, catalog)
	}
}

func GetCatalogFiles(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.CatalogsResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		files := getCatalogChildFilesOr404(db, id, w, r)
		if files == nil {
			return
		}

		respondJSON(w, http.StatusOK, files)
	}
}

func GetCatalogCatalogs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.CatalogsResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		catalogs := getCatalogChildCatalogsOr404(db, id, w, r)
		if catalogs == nil {
			return
		}

		respondJSON(w, http.StatusOK, catalogs)
	}
}

func GetCatalogPath(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var path []model.CatalogPathElement

	vars := mux.Vars(r)
	hashId := vars["hashId"]
	zeroId := utils.EncodeId(0, utils.CatalogsResourceType)
	index := 0

	if hashId == zeroId {
		respondError(w, http.StatusBadRequest, "RootHasNoPath")
		return
	}

	for len(path) == 0 || hashId != zeroId {
		id, err := utils.DecodeId(hashId, utils.CatalogsResourceType)

		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		} else {
			var catalog *model.Catalog

			catalog = getCatalogOr404(db, uint(id), w, r)
			if catalog == nil {
				return
			}

			catalog = catalog.WithHashId()
			path = append(path, catalog.ToPath(index))
			hashId = catalog.ParentHashID
			index++
		}
	}

	// last parent is always root (titled as 'Files')
	path = append(path, model.CatalogPathElement{
		Index:    index + 1,
		HashId:   zeroId,
		Title:    "Files",
		IsPublic: true,
	})

	respondJSON(w, http.StatusOK, path)
}

func getRootCatalogsOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) []model.Catalog {
	var catalogs []model.Catalog
	if err := db.
		Where("is_public = 1 AND parent_id = 0").
		Find(&catalogs, &model.Catalog{}).
		Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	for i, _ := range catalogs {
		catalogs[i].WithHashId()
	}
	return catalogs
}

func getCatalogOr404(db *gorm.DB, id uint, w http.ResponseWriter, r *http.Request) *model.Catalog {
	var catalog model.Catalog
	if err := db.
		First(&catalog, model.Catalog{
			Model: model.Model{ID: id},
		}).
		Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return catalog.WithHashId()
}

func getCatalogChildCatalogsOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) []model.Catalog {
	var catalogs []model.Catalog
	if err := db.
		Where(model.Catalog{ParentID: id}).
		Find(&catalogs, model.Catalog{}).
		Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	for i, _ := range catalogs {
		catalogs[i].WithHashId()
	}
	return catalogs
}

func getCatalogChildFilesOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) []model.File {
	var files []model.File
	if err := db.
		Where(model.File{CatalogID: id}).
		Find(&files, model.File{}).
		Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	for i, _ := range files {
		files[i].WithHashId().WithHasPassword()
	}
	return files
}
