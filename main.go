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
	r.Static("/uploads", "./uploads")
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshToken)
	r.GET("/", middlewares.CheckAuth, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted to protected route"})
	})
	runErr := r.Run()
	if runErr != nil {
		return
	}
}
