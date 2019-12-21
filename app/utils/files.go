package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/AdrianOrlow/files-api/config"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func InitializeFiles(config *config.Config) error {
	utils.files.dir = config.FilesDir
	err := createFilesDirIfNotExists()

	return err
}

func SaveMultipartFile(file multipart.File, fileName string, w http.ResponseWriter, r *http.Request) error {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(utils.files.dir+"/"+fileName, fileData, 0666)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(fileName string) error {
	err := os.Remove(utils.files.dir + "/" + fileName)
	return err
}

func ReadFile(fileName string) (*os.File, error) {
	data, err := os.Open(utils.files.dir + "/" + fileName)
	return data, err
}

func ReadFileByteStream(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(utils.files.dir + "/" + fileName)
	return data, err
}

func GetFileNameWithTimestamp(fileName string) string {
	extension := filepath.Ext(fileName)
	name := strings.TrimRight(fileName, extension)
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	return name + "-" + timestamp + extension
}

func GetFileNameWithoutTimestamp(fileName string) string {
	re := regexp.MustCompile("(.+)(-[0-9]+)(.[a-z]+)")

	split := re.FindAllStringSubmatch(fileName, -1)
	return split[0][1] + split[0][3]
}

func GetFileMD5Hash(file io.Reader) (string, error) {
	md5Hash := md5.New()
	_, err := io.Copy(md5Hash, file)
	if err != nil {
		return "", err
	}
	fileMd5HashInBytes := md5Hash.Sum(nil)[:16]
	return hex.EncodeToString(fileMd5HashInBytes), nil
}

func GetFileSHA1Hash(file io.Reader) (string, error) {
	sha1Hash := sha1.New()
	_, err := io.Copy(sha1Hash, file)
	if err != nil {
		return "", err
	}
	fileSha1HashInBytes := sha1Hash.Sum(nil)[:20]
	return hex.EncodeToString(fileSha1HashInBytes), nil
}

func GetFileSizeInKilobytes(file *os.File) (string, error) {
	savedFileStat, err := file.Stat()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(savedFileStat.Size()/1024, 10), nil
}

func createFilesDirIfNotExists() error {
	_, err := os.Stat(utils.files.dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(utils.files.dir, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
