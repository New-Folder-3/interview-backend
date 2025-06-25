package middlewares

import (
	"github.com/gin-gonic/gin"
	"interview-backend/util"
	"io"
	"net/http"
	"strings"
)

func APINoRoute(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}

func FileChecker(typ string) func(c *gin.Context) {
	return func(c *gin.Context) {
		switch typ {
		case "image":
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20)
		case "audio":
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)
		}
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 400, false)
			c.Abort()
			return
		}
		defer file.Close()
		buf := make([]byte, 512)
		_, err = file.Read(buf)
		if err != nil && err != io.EOF {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 400, false)
			c.Abort()
		}
		mimeType := http.DetectContentType(buf)
		switch typ {
		case "image":
			if mimeType != "image/jpeg" && mimeType != "image/png" {
				util.ErrorResp(c, "Not a valid image type", 400, false)
				c.Abort()
				return
			}
		case "audio":
			if mimeType != "audio/mpeg" && mimeType != "audio/wav" {
				util.ErrorResp(c, "Not a valid audio type", 400, false)
				c.Abort()
				return
			}
		}
		file.Seek(0, io.SeekStart)
		c.Set("reFile", file)
		c.Set("reHeader", header)
		c.Next()
	}
}
