package file

import (
	"fmt"
	"os"
	"strings"
)
var (
	UploadFileDir string = "/uploadfile"
	FileServerUrl = "/file/download_file"
)

type UploadFileResponse struct {
	FileUrls []string `json:"file_urls"`
}

func GetUploadFileName(sz int64, oldName string) string {
	fn := fmt.Sprintf("%d-%s", sz, oldName)
	fn = strings.ReplaceAll(fn, "/", ".")
	fn = strings.ReplaceAll(fn, "\\", ".")
	fn = strings.ReplaceAll(fn, " ", ".")
	fn = strings.ReplaceAll(fn, "(", ".")
	fn = strings.ReplaceAll(fn, ")", ".")
	return fn
	//return util.RandomString(48)
}

func GetDownloadFileUrl(fileName string) string {
	return strings.Join([]string{FileServerUrl,"?", fileName},"")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
