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
			fmt.Println(err) // Logging the error for debugging purposes
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
