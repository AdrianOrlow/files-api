package handler

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/AdrianOrlow/files-api/app/model"
	"github.com/AdrianOrlow/files-api/app/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetFile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FilesResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		file := getFileOr404(db, uint(id), w, r)
		if file == nil {
			return
		}

		respondJSON(
			w,
			http.StatusOK,
			file.WithHashId().WithHasPassword().WithFileNameWithoutTimestamp(),
		)
	}
}

func CreateFile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var file model.File

	err := r.ParseMultipartForm(128 * 8 << 20) // 128 MB
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	multipartFile, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal([]byte(r.FormValue("data")), &file)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.Save(&file).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	file.FileName = utils.GetFileNameWithTimestamp(file.FileName)
	err = utils.SaveMultipartFile(multipartFile, file.FileName, w, r)
	if err != nil {
		deleteErr := db.Delete(&file).Error
		if deleteErr != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	savedFile, err := utils.ReadFile(file.FileName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	file.FileChecksumSHA1, err = utils.GetFileSHA1Hash(savedFile)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	savedFile, err = utils.ReadFile(file.FileName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	file.FileChecksumMD5, err = utils.GetFileMD5Hash(savedFile)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	file.FileSizeKB, err = utils.GetFileSizeInKilobytes(savedFile)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var password model.FilePassword
	err = json.Unmarshal([]byte(r.FormValue("data")), &password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if password.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.MinCost)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		file.Password = string(hashPassword)
	}

	err = db.Save(file.WithId()).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	savedFile.Close()
	multipartFile.Close()
	respondJSON(
		w,
		http.StatusCreated,
		file.WithHasPassword().WithHashId().WithFileNameWithoutTimestamp(),
	)
}

func DeleteFile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FilesResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		file := getFileOr404(db, uint(id), w, r)
		if file == nil {
			return
		}

		err := utils.DeleteFile(file.FileName)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = db.Delete(&file).Error
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusNoContent, nil)
	}
}

func ServeFile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FilesResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		file := getFileOr404(db, uint(id), w, r)
		if file == nil {
			return
		}

		if file.WithHasPassword().HasPassword {
			key, ok := vars["key"]
			if !ok {
				respondError(w, http.StatusUnauthorized, "key not provided")
				return
			}

			link := getFileLinkOr401(db, key, w, r)
			if link == nil {
				return
			}

			if file.ID != uint(link.FileID) {
				respondError(w, http.StatusNotFound, "wrong file id")
				return
			}

			if !link.IsValid() {
				respondError(w, http.StatusNotFound, "link expired")
				return
			}
		}

		data, err := utils.ReadFileByteStream(file.FileName)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		file = file.WithFileNameWithoutTimestamp()
		w.Header().Set("Content-Disposition", "attachment; filename=" + file.FileName)
		http.ServeContent(w, r, file.FileName, time.Now(), bytes.NewReader(data))
	}
}

func CreateFileDownloadLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hashId := vars["hashId"]
	id, err := utils.DecodeId(hashId, utils.FilesResourceType)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		file := getFileOr404(db, uint(id), w, r)
		if file == nil {
			return
		}

		if file.WithHasPassword().HasPassword {
			err := utils.CompareHashAndPasswordFromAuthHeader([]byte(file.Password), r)
			if err != nil {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		key := utils.GenerateRandomString(25)
		link := model.Link{
			Key: key,
			FileID: int(file.ID),
			ValidUntil: time.Now().Add(time.Hour * time.Duration(100)),
		}

		err = db.Save(&link).Error
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, link.WithHashId())
	}
}

func getFileOr404(db *gorm.DB, id uint, w http.ResponseWriter, r *http.Request) *model.File {
	var file model.File
	err := db.
		First(&file, model.File{Model: model.Model{ID: id}}).
		Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return file.WithHashId()
}

func getFileLinkOr401(db *gorm.DB, key string, w http.ResponseWriter, r *http.Request) *model.Link {
	var link model.Link
	err := db.
		First(&link, model.Link{Key: key}).
		Error
	if err != nil {
		respondError(w, http.StatusUnauthorized, "valid key not provided")
		return nil
	}
	return link.WithHashId()
}
