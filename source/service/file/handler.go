package file

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// "strings"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func GetFile(c *gin.Context) {
	//这样子验证有问题，还是任何用户都可以获取到.
	FileName := c.Query("file_name")
	FileName = strings.ReplaceAll(FileName,"/",".")
	FileName = strings.ReplaceAll(FileName,"\\",".")
	FilePath := filepath.Join(UploadFileDir, FileName)
	fi, err := os.Stat(FilePath)
	if err != nil || fi.IsDir() || fi.Size() == 0 {
		c.JSON(http.StatusBadRequest,"file path is error !")
		return;
	}
	c.File(FilePath)
}

func UploadFile(c *gin.Context) {
	file, e := c.FormFile("file")
	uploadFileResponse := &UploadFileResponse{}
	if e != nil {
		log.Errorf(e, "Has A Error in FormFile")
		c.JSON(http.StatusBadRequest, uploadFileResponse)
		return
	}
	os.Mkdir(UploadFileDir, 0777)
	filename := GetUploadFileName(file.Size, file.Filename)
	filePath := filepath.Join(UploadFileDir, filename)
	log.Infof("upload filename old:%s random name:%s", file.Filename, filename)
	fileUrls := make([]string, 1)
	if bExist, _ := PathExists(filePath); !bExist {
		if e = c.SaveUploadedFile(file, filePath); e != nil {
			log.Errorf(e, "Has A Error in SaveUploadedFile")
			c.JSON(http.StatusBadRequest, uploadFileResponse)
			return
		}
	} else {
		log.Debugf("filePath:%s is exist alaready .", filePath)
		fileUrls[0] = GetDownloadFileUrl(filename)
		uploadFileResponse.FileUrls = fileUrls
		c.AbortWithStatusJSON(http.StatusOK, uploadFileResponse)
		return
	}
	log.Debugf("upload file success file:%s", filePath)
	fileUrls[0] =  GetDownloadFileUrl(filename)
	uploadFileResponse.FileUrls = fileUrls

	c.JSON(http.StatusOK, uploadFileResponse)
}


func UploadFileList(c *gin.Context) {
	file, e := c.FormFile("file")
	uploadFileResponse := &UploadFileResponse{}
	if e != nil {
		log.Errorf(e, "Has A Error in FormFile")
		c.JSON(http.StatusBadRequest, uploadFileResponse)
		return
	}
	os.Mkdir(UploadFileDir, 0777)
	filename := GetUploadFileName(file.Size, file.Filename)
	filePath := filepath.Join(UploadFileDir, filename)
	log.Infof("upload filename old:%s random name:%s", file.Filename, filename)
	fileUrls := make([]string, 1)
	if bExist, _ := PathExists(filePath); !bExist {
		if e = c.SaveUploadedFile(file, filePath); e != nil {
			log.Errorf(e, "Has A Error in SaveUploadedFile")
			c.JSON(http.StatusBadRequest, uploadFileResponse)
			return
		}
	} else {
		log.Debugf("filePath:%s is exist alaready .", filePath)
		fileUrls[0] = GetDownloadFileUrl(filename)
		uploadFileResponse.FileUrls = fileUrls
		c.AbortWithStatusJSON(http.StatusOK, uploadFileResponse)
		return
	}
	log.Debugf("upload file success file:%s", filePath)
	fileUrls[0] =  GetDownloadFileUrl(filename)
	uploadFileResponse.FileUrls = fileUrls

	c.JSON(http.StatusOK, uploadFileResponse)
}
