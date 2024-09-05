package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/initializers"
	"main/models"
	"net/http"
	"strconv"
)

func FollowUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID cannot be empty"})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
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

	var followingUser models.User

	err = initializers.DB.Where("id = ?", intID).First(&followingUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	var followExist models.FollowModel

	errFollow := initializers.DB.Where("followed_id = ? AND following_id = ?", userModel.ID, followingUser.ID).First(&followExist).Error
	if errFollow == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already following this user"})
		return
	}

	if followingUser.ID == userModel.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself!"})
		return
	}

	follow := models.FollowModel{
		FollowingID:  followingUser.ID,
		FollowedByID: userModel.ID,
	}

	err = initializers.DB.Create(&follow).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}

func UnFollow(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID cannot be empty"})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
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

	var followingUser models.User
	err = initializers.DB.Where("id = ?", intID).First(&followingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	var followExist models.FollowModel
	errFollow := initializers.DB.Where("followed_by_id = ? AND following_id = ?", followingUser.ID, userModel.ID).First(&followExist).Error
	if errFollow == nil {
		if err = initializers.DB.Delete(&followExist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
			return
		}
	} else if errors.Is(errFollow, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not following this user"})
		return
	} else {
		fmt.Println(errFollow)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Unfollowed successfully"})
}

func ListFollowers(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	currentUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var followers []models.User
	err := initializers.DB.Model(&models.FollowModel{}).
		Where("followed_by_id = ?", currentUser.ID).
		Joins("JOIN users ON users.id = follow_models.following_id").
		Find(&followers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve followers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
	})
}

func ListFollowings(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	currentUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var followings []models.User
	err := initializers.DB.Model(&models.FollowModel{}).
		Where("following_id = ?", currentUser.ID).
		Joins("JOIN users ON users.id = follow_models.followed_by_id").
		Find(&followings).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve followings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followings": followings,
	})
}

func LikeTweet(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID cannot be empty"})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
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
	err = initializers.DB.Where("id = ?", intID).First(&tweet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	var like models.LikeModel
	err = initializers.DB.Where("user_id = ? AND tweet_id = ?", userModel.ID, tweet.ID).First(&like).Error

	if err == nil {
		// Like already exists, so return an error
		c.JSON(http.StatusConflict, gin.H{"error": "Tweet already liked"})
		return
	}

	// Create the like
	like = models.LikeModel{
		UserID:  userModel.ID,
		TweetID: tweet.ID,
	}

	if err := initializers.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like tweet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tweet liked successfully"})
}

func UnlikeTweet(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID cannot be empty"})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
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
	err = initializers.DB.Where("id = ?", intID).First(&tweet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		}
		return
	}

	var like models.LikeModel
	err = initializers.DB.Where("user_id = ? AND tweet_id = ?", userModel.ID, tweet.ID).First(&like).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Like does not exist, so return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	// Delete the like
	if err := initializers.DB.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike tweet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tweet unliked successfully"})
}
