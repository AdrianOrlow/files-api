package handler

import (
	"encoding/json"
	"github.com/AdrianOrlow/files-api/app/utils"
	"net/http"

	"github.com/AdrianOrlow/files-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetRootFolders(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	folders := getRootFoldersOr404(db, w, r)
	if folders == nil {
		return
	}
	respondJSON(w, http.StatusOK, folders)
}

func CreateFolder(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var folder model.Folder

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&folder)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = db.Save(folder.WithId()).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, folder.WithHashId())
}

func UpdateFolder(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FoldersResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		folder := getFolderOr404(db, uint(id), w, r)
		if folder == nil {
			return
		}

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&folder)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		err = db.Save(folder.WithId()).Error
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, folder.WithHashId())
	}
}

func DeleteFolder(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FoldersResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		folder := getFolderOr404(db, uint(id), w, r)
		if folder == nil {
			return
		}
		err = db.Delete(&folder).Error
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusNoContent, nil)
	}
}

func GetFolder(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FoldersResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		folder := getFolderOr404(db, uint(id), w, r)
		if folder == nil {
			return
		}
		respondJSON(w, http.StatusOK, folder)
	}
}

func GetFolderFiles(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FoldersResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		files := getFolderChildFilesOr404(db, id, w, r)
		if files == nil {
			return
		}

		respondJSON(w, http.StatusOK, files)
	}
}

func GetFolderFolders(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FoldersResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		folders := getFolderChildFoldersOr404(db, id, w, r)
		if folders == nil {
			return
		}

		respondJSON(w, http.StatusOK, folders)
	}
}

func GetFolderPath(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var path []model.FolderPathElement

	vars := mux.Vars(r)
	hashId := vars["hashId"]
	zeroId := utils.EncodeId(0, utils.FoldersResourceType)
	index := 0

	if hashId == zeroId {
		respondError(w, http.StatusBadRequest, "RootHasNoPath")
		return
	}

	for len(path) == 0 || hashId != zeroId {
		id, err := utils.DecodeId(hashId, utils.FoldersResourceType)

		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		} else {
			var folder *model.Folder

			folder = getFolderOr404(db, uint(id), w, r)
			if folder == nil {
				return
			}

			folder = folder.WithHashId()
			path = append(path, folder.ToPath(index))
			hashId = folder.ParentHashID
			index++
		}
	}

	respondJSON(w, http.StatusOK, path)
}

func getRootFoldersOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) []model.Folder {
	var folders []model.Folder
	err := db.
		Where("is_public = 1 AND parent_id = 0").
		Find(&folders, &model.Folder{}).
		Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	for i, _ := range folders {
		folders[i].WithHashId()
	}
	return folders
}

func getFolderOr404(db *gorm.DB, id uint, w http.ResponseWriter, r *http.Request) *model.Folder {
	var folder model.Folder
	err := db.
		First(&folder, model.Folder{
			Model: model.Model{ID: id},
		}).
		Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return folder.WithHashId()
}

func getFolderChildFoldersOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) []model.Folder {
	var folders []model.Folder
	err := db.
		Where(model.Folder{ParentID: id}).
		Find(&folders, model.Folder{}).
		Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	for i, _ := range folders {
		folders[i].WithHashId()
	}
	return folders
}

func getFolderChildFilesOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) []model.File {
	var files []model.File
	err := db.
		Where(model.File{FolderID: id}).
		Find(&files, model.File{}).
		Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	for i, _ := range files {
		files[i].WithHashId().WithHasPassword()
	}
	return files
}
