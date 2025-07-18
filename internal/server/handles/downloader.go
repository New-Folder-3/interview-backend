package handles

import (
	"github.com/gin-gonic/gin"
	"interview/cmd/flags"
	"interview/internal/conf"
	"interview/util"
	"io"
	"path"
	"path/filepath"
	"strings"
)

func DownloadHandler(c *gin.Context) {
	reqPath := strings.TrimPrefix(c.Request.URL.Path, "/api/download/")
	if conf.Conf.S3.Use {
		filePath := path.Join(conf.Conf.S3.Path, reqPath)
		reader, err := util.S3Download(filePath)
		if err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, "S3 download error", 404)
			return
		}
		defer reader.Close()

		c.Header("Content-Disposition", "attachment; filename=\""+path.Base(reqPath)+"\"")
		c.Header("Content-Type", "application/octet-stream")

		_, err = io.Copy(c.Writer, reader)
		if err != nil {
			util.ErrorPrinter(err)
			return
		}
	} else {
		filePath := path.Join(flags.DataDir, reqPath)
		if _, err := filepath.Abs(filePath); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, "Invalid file path", 404)
			return
		}
		c.File(filePath)
	}
}
