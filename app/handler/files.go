package handler

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		file := getFileOr404(db, uint(id), w, r).WithHasPassword()
		if file == nil {
			return
		}

		if file.HasPassword {
			err := checkPasswordFromAuthorizationHeader(file.Password, w, r)
			if err != nil {
				return
			}
		}

		respondJSON(w, http.StatusOK, file.WithHashId().WithHasPassword())
	}
}

func CreateFile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var file model.File

	err := r.ParseMultipartForm(128 << 20) // 128 Mb
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	reqFile, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer reqFile.Close()

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

	file.FileName = getFileNameWithTimestamp(file.FileName)
	err = saveFile(reqFile, file.FileName, w, r)
	if err != nil {
		err := db.Delete(&file).Error
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	savedFile, err := os.Open("./files/" + file.FileName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer savedFile.Close()
	file.FileChecksumSHA1, err = getFileSHA1Hash(savedFile)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	savedFile, _ = os.Open("./files/" + file.FileName)
	file.FileChecksumMD5, err = getFileMD5Hash(savedFile)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	savedFileStat, err := savedFile.Stat()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	file.FileSizeKB = strconv.FormatInt(savedFileStat.Size()/1024, 10)

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

	err = db.Save(&file).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, file.WithHasPassword().WithHashId())
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

		err := os.Remove("./files/" + file.FileName)
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
		file := getFileOr404(db, uint(id), w, r).WithHasPassword()
		if file == nil {
			return
		}

		if file.HasPassword {
			err := checkPasswordFromAuthorizationHeader(file.Password, w, r)
			if err != nil {
				return
			}
		}

		data, err := ioutil.ReadFile("files/" + file.FileName)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		http.ServeContent(w, r, file.FileName, time.Now(), bytes.NewReader(data))
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

func saveFile(file multipart.File, fileName string, w http.ResponseWriter, r *http.Request) error {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return err
	}

	err = ioutil.WriteFile("./files/"+fileName, fileData, 0666)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

func getFileNameWithTimestamp(fileName string) string {
	extenstion := filepath.Ext(fileName)
	name := strings.TrimRight(fileName, extenstion)
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	return name + "-" + timestamp + extenstion
}

func getFileMD5Hash(file io.Reader) (string, error) {
	md5Hash := md5.New()
	_, err := io.Copy(md5Hash, file)
	if err != nil {
		return "", err
	}
	fileMd5HashInBytes := md5Hash.Sum(nil)[:16]
	return hex.EncodeToString(fileMd5HashInBytes), nil
}

func getFileSHA1Hash(file io.Reader) (string, error) {
	sha1Hash := sha1.New()
	_, err := io.Copy(sha1Hash, file)
	if err != nil {
		return "", err
	}
	fileSha1HashInBytes := sha1Hash.Sum(nil)[:20]
	return hex.EncodeToString(fileSha1HashInBytes), nil
}

func checkPasswordFromAuthorizationHeader(filePassword string, w http.ResponseWriter, r *http.Request) error {
	password := getPasswordFromAuthorizationHeader(r)
	if password != nil {
		err := bcrypt.CompareHashAndPassword([]byte(filePassword), password)

		if err != nil {
			respondError(w, http.StatusForbidden, "PasswordIncorrect")
			return err
		}
	} else {
		respondError(w, http.StatusBadRequest, "AuthorizationHeaderNotProvided")
		return errors.New("AuthorizationHeaderNotProvided")
	}
	return nil
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
