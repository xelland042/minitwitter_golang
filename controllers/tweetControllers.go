package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"main/initializers"
	"main/models"
	"main/utils"
	"net/http"
	"path/filepath"
	"time"
)

func CreateTweet(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	Title := c.Request.FormValue("title")
	Body := c.Request.FormValue("body")

	var filePath string
	file, err := c.FormFile("file")
	if errors.Is(err, http.ErrMissingFile) {
		filePath = ""
	}

	if file != nil {
		filePath = utils.GetUniqueFileName("uploads/tweets/", time.Now().Format("20060102150405"), filepath.Ext(file.Filename))
		errFile := c.SaveUploadedFile(file, filePath)
		if errFile != nil {
			return
		}
	} else {
		filePath = ""
	}

	tweet := models.Tweet{
		Title: Title,
		Body:  Body,
		File:  filePath,
	}

	initializers.DB.Create(&tweet)

	c.JSON(http.StatusOK, gin.H{"tweet": tweet})
}
