package server

import (
	"context"

	controllers "github.com/MoulieshN/Go-JWT-Project.git/controllers"
	"github.com/MoulieshN/Go-JWT-Project.git/repository"
	"github.com/gin-gonic/gin"
)

func NewRoutes(c context.Context, repo repository.UserRepository) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	// User-related routes
	authorized := router.Group("/api/v1/auth")
	UserController := controllers.NewUserController(repo)

	authorized.POST("user/signup", UserController.SignUp())
	authorized.POST("user/login", UserController.Login())

	// Add authentication middleware only to internal routes

	internal := router.Group("/api/v1/users")
	// router.Use(middleware.Authenticate())
	internal.GET("", UserController.GetUsers())
	internal.GET("/:id", UserController.GetUser())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	return router
}
