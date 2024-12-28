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

	authorized := router.Group("/api/v1/auth")
	internal := router.Group("/api/v1/users")

	UserController := controllers.NewUserController(repo)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	authorized.POST("user/signup", controller.SignUp)
	authorized.POST("user/login", controller.Login)

	router.Use(middleware.Authenticate())
	internal.GET("", UserController.GetUsers())
	internal.GET("/:id", UserController.GetUser())

	return router
}
