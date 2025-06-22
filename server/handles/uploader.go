package handles

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"interview-backend/cmd/flags"
	"interview-backend/internal/conf"
	"net/http"
	"path/filepath"
)

func PicUploaderHandle(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file is received"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only jpg/jpeg/png images are allowed"})
		return
	}

	filename := fmt.Sprintf("%d%s", uuid.New().String(), ext)
	filePath := filepath.Join(flags.DataDir, "image", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	URL := fmt.Sprintf("%s/api/image/%s", conf.Conf.Schema.URL, filename)
	c.JSON(http.StatusOK, gin.H{
		"url": URL,
	})
}

func AudioUploaderHandle(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No audio file is received"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".mp3" && ext != ".wav" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only mp3/wav images are allowed"})
		return
	}

	filename := fmt.Sprintf("%d%s", uuid.New().String(), ext)
	filePath := filepath.Join(flags.DataDir, "audio", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	URL := fmt.Sprintf("%s/api/audio/%s", conf.Conf.Schema.URL, filename)
	c.JSON(http.StatusOK, gin.H{
		"url": URL,
	})
}
