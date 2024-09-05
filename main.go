package main

import (
	"github.com/gin-gonic/gin"
	"main/controllers"
	"main/initializers"
	"main/middlewares"
	"net/http"
)

func init() {
	initializers.LoadEnVVariables()
	initializers.ConnectToDB()
	initializers.SyncDataBase()
}

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 30 << 20
	uploads := r.Group("/uploads")
	uploads.Use(middlewares.CheckAuth) // Apply your authentication middleware
	{
		uploads.GET("/*filepath", func(c *gin.Context) {
			filepath := c.Param("filepath")
			c.File("./uploads" + filepath)
		})
	}
	//Users endpoints
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshToken)
	r.GET("/user", middlewares.CheckAuth, controllers.UserProfile)
	r.POST("/user", middlewares.CheckAuth, controllers.UserUpdate)
	r.POST("/change-password", middlewares.CheckAuth, controllers.ChangePassword)
	r.GET("/", middlewares.CheckAuth, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted to protected route"})
	})

	// Tweets endpoint
	r.POST("/create-tweet", middlewares.CheckAuth, controllers.CreateTweet)
	runErr := r.Run()
	if runErr != nil {
		return
	}
}
