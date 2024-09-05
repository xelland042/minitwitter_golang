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
	r.MaxMultipartMemory = 10 << 20
	authorized := r.Group("/uploads")
	authorized.Use(middlewares.CheckAuth) // Apply your authentication middleware
	{
		authorized.GET("/*filepath", func(c *gin.Context) {
			filepath := c.Param("filepath")
			c.File("./uploads" + filepath)
		})
	}
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshToken)
	r.GET("/user", middlewares.CheckAuth, controllers.UserProfile)
	r.GET("/", middlewares.CheckAuth, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted to protected route"})
	})
	runErr := r.Run()
	if runErr != nil {
		return
	}
}
