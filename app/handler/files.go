package handler

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"

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
		file := getFileOr404(db, uint(id), w, r).WithHasPassword()
		if file == nil {
			return
		}

		if file.HasPassword {
			password := getPasswordFromAuthorizationHeader(r)
			if password != nil {
				err := bcrypt.CompareHashAndPassword([]byte(file.Password), password)

				if err != nil {
					respondError(w, http.StatusForbidden, "PasswordIncorrect")
					return
				}
			} else {
				respondError(w, http.StatusBadRequest, "AuthorizationHeaderNotProvided")
				return
			}
		}

		respondJSON(w, http.StatusOK, file.WithHashId().WithHasPassword())
	}
}

func getFileOr404(db *gorm.DB, id uint, w http.ResponseWriter, r *http.Request) *model.File {
	var file model.File
	if err := db.
		First(&file, model.File{Model:
			model.Model{ID: id},
		}).
		Error;
	err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return file.WithHashId()
}

func getPasswordFromAuthorizationHeader(r *http.Request) []byte {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return nil
	}
	splitToken := strings.Split(reqToken, "Basic ")
	reqToken = splitToken[1]
	password, _ := base64.StdEncoding.DecodeString(reqToken)
	return password
}
