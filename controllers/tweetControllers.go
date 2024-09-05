package controllers

import "C"
import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
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
		Title:    Title,
		Body:     Body,
		File:     filePath,
		AuthorID: userModel.ID,
	}

	initializers.DB.Create(&tweet)

	response := utils.TweetResponse{
		ID:        tweet.ID,
		CreatedAt: tweet.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tweet.UpdatedAt.Format(time.RFC3339),
		Title:     tweet.Title,
		Body:      tweet.Body,
		File:      tweet.File,
	}

	c.JSON(http.StatusOK, gin.H{"tweet": response})
}

func TweetList(c *gin.Context) {
	searchQuery := c.Query("search")

	var tweets []map[string]interface{}

	query := initializers.DB.Model(&models.Tweet{}).Select("id, title, body, created_at")

	if searchQuery != "" {
		query = query.Where("title ILIKE ? OR body ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	err := query.Find(&tweets).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets": tweets,
	})
}

func TweetRetrieve(c *gin.Context) {
	id := c.Param("id")

	var tweet models.Tweet

	err := initializers.DB.Where("id = ?", id).First(&tweet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	response := utils.TweetResponse{
		ID:        tweet.ID,
		CreatedAt: tweet.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tweet.UpdatedAt.Format(time.RFC3339),
		Title:     tweet.Title,
		Body:      tweet.Body,
		File:      tweet.File,
	}

	c.JSON(http.StatusOK, gin.H{
		"tweet": response,
	})
}

func TweetUpdate(c *gin.Context) {
	id := c.Param("id")

	if err := c.Request.ParseMultipartForm(30 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var tweet models.Tweet

	err := initializers.DB.Where("id = ? AND author_id = ?", id, userModel.ID).First(&tweet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	title := c.Request.FormValue("title")
	body := c.Request.FormValue("body")

	var filePath string
	file, err := c.FormFile("file")
	if err == nil {
		fileName := time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
		filePath = filepath.Join("uploads/tweets", fileName)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	} else {
		filePath = tweet.File
	}

	fmt.Println(filePath)

	if title != "" {
		tweet.Title = title
	}
	if body != "" {
		tweet.Body = body
	}
	if filePath != "" {
		tweet.File = filePath
	}

	if err := initializers.DB.Save(&tweet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"title": tweet.Title,
			"body":  tweet.Body,
			"file":  tweet.File,
		},
	})
}

func TweetDelete(c *gin.Context) {
	id := c.Param("id")

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var tweet models.Tweet

	err := initializers.DB.Where("id = ? AND author_id = ?", id, userModel.ID).First(&tweet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	if err := initializers.DB.Delete(&tweet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tweet"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Tweet deleted successfully"})
}
