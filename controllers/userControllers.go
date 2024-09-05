package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"main/initializers"
	"main/models"
	"main/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SignUp(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	UserName := c.Request.FormValue("UserName")
	Password := c.Request.FormValue("Password")
	Email := c.Request.FormValue("Email")
	Bio := c.Request.FormValue("Bio")

	var userFound models.User
	initializers.DB.Where("username=?", UserName).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already used"})
		return
	}

	if !utils.IsValidEmail(Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if len(Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	if !utils.IsStrongPassword(Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"})
		return
	}

	passwordHash, errPassword := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if errPassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errPassword.Error()})
		return
	}

	var filePath string
	file, err := c.FormFile("Picture")
	if errors.Is(err, http.ErrMissingFile) {
		filePath = ""
	}

	if file != nil {
		filePath = utils.GetUniqueFileName("uploads/profile_pictures/", time.Now().Format("20060102150405"), filepath.Ext(file.Filename))
		errFile := c.SaveUploadedFile(file, filePath)
		if errFile != nil {
			return
		}
	} else {
		filePath = ""
	}

	user := models.User{
		UserName: UserName,
		Email:    Email,
		Password: string(passwordHash),
		Bio:      Bio,
		Picture:  filePath,
	}

	initializers.DB.Create(&user)

	response := utils.UserResponse{
		UserName: user.UserName,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK,
		gin.H{"data": response})
}

func Login(c *gin.Context) {
	var loginInput utils.LoginInput

	if errAuthIn := c.ShouldBindJSON(&loginInput); errAuthIn != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errAuthIn.Error()})
		return
	}

	var user models.User
	var errLogin error

	if loginInput.Email != "" {
		errLogin = initializers.DB.Where("email = ?", loginInput.Email).First(&user).Error
	} else if loginInput.UserName != "" {
		errLogin = initializers.DB.Where("username = ?", loginInput.UserName).First(&user).Error
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either username or email must be provided"})
		return
	}

	if errLogin != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); errPassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	accessTokenClaims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshTokenClaims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", accessTokenString, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})
}

func RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userId := claims["id"].(float64)

	accessTokenClaims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessTokenString,
	})
}

func UserProfile(c *gin.Context) {
	user, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	u, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User data is invalid"})
		return
	}

	profilePictureURL := ""
	if u.Picture != "" {
		profilePictureURL = u.Picture
	}

	c.JSON(http.StatusOK, gin.H{
		"username": u.UserName,
		"email":    u.Email,
		"bio":      u.Bio,
		"picture":  profilePictureURL,
	})
}

func UserUpdate(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(models.User)

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	username := c.Request.FormValue("UserName")
	email := c.Request.FormValue("Email")
	bio := c.Request.FormValue("Bio")

	var filePath string
	file, err := c.FormFile("Picture")
	if err == nil {
		fileName := time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
		filePath = filepath.Join("uploads", fileName)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	} else {
		filePath = currentUser.Picture
	}

	fmt.Println(username)

	if username != "" {
		currentUser.UserName = username
	}
	if email != "" {
		currentUser.Email = email
	}
	if bio != "" {
		currentUser.Bio = bio
	}
	if filePath != "" {
		currentUser.Picture = filePath
	}

	if err := initializers.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"username": currentUser.UserName,
			"email":    currentUser.Email,
			"bio":      currentUser.Bio,
			"picture":  currentUser.Picture,
		},
	})
}
