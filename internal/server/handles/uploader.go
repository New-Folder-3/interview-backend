package handles

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"interview/cmd/flags"
	"interview/internal/conf"
	"interview/util"
	"mime/multipart"
	"net/http"
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
		file := fileRe.(multipart.File)
		header := headerRe.(*multipart.FileHeader)

		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(flags.DataDir, typ, filename)

		err := util.SaveUploadFile(filePath, &file)
		if err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
		}

		URL := conf.Conf.Schema.URL + path.Join("/api/download", typ, filename)
		c.JSON(http.StatusOK, gin.H{
			"url": URL,
		})
	}
}
