package handles

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"interview/cmd/flags"
	"interview/internal/conf"
	"interview/util"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func FileSaver(typ string) func(*gin.Context) {
	return func(c *gin.Context) {
		fileRe, ok1 := c.Get("reFile")
		headerRe, ok2 := c.Get("reHeader")
		if !ok1 || !ok2 {
			util.ErrorResp(c, "Incomplete upload", 500)
			return
		}
		file := fileRe.(io.Reader)
		header := headerRe.(*multipart.FileHeader)

		ext := filepath.Ext(header.Filename)
		if typ == "video" {
			var err error
			file, err = util.WebmToMp4Stream(file)
			if err != nil {
				util.ErrorPrinter(err)
				util.ErrorResp(c, "Failed to convert video format", 500)
				return
			}
			ext = ".mp4"
		}

		var err error
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		tmpfile, err := os.CreateTemp("", "upload-*")
		if err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, "Failed to create temp file", 500)
			return
		}
		defer os.Remove(tmpfile.Name())
		defer tmpfile.Close()

		if _, err = io.Copy(tmpfile, file); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, "Failed to save file", 500)
			return
		}
		if _, err = tmpfile.Seek(0, io.SeekStart); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, "Failed to seek file", 500)
			return
		}

		if conf.Conf.S3.Use {
			filePath := filepath.Join(conf.Conf.S3.Path, filename)
			err = util.S3Upload(tmpfile, filePath)
		} else {
			filePath := filepath.Join(flags.DataDir, typ, filename)
			err = util.SaveUploadFile(filePath, tmpfile)
		}

		if err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
			return
		}

		URL := conf.Conf.Schema.URL + path.Join("/api/download", typ, filename)
		util.SuccessResp(c, URL, "File saved successfully")
	}
}

func Audio64Saver(c *gin.Context) {
	var request UploadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	fileName := uuid.New().String() + ".mp3"
	filePath := filepath.Join(flags.DataDir, "audio", fileName)
	err := util.SaveBase64WebmToMp3(request.Base64, filePath)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	URL := conf.Conf.Schema.URL + path.Join("/api/download", "audio", fileName)
	util.SuccessResp(c, URL, "Audio saved successfully")
}
