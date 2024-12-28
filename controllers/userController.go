package controllers

import (
	"net/http"

	"github.com/MoulieshN/Go-JWT-Project.git/helpers"
	"github.com/MoulieshN/Go-JWT-Project.git/models"
	"github.com/MoulieshN/Go-JWT-Project.git/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) UserController {
	return UserController{userRepo: repo}
}

func HashPassword() {
	return
}

func VerfiyPassword() {
	return
}

func (u *UserController) SignUp() {
	return
}

func (u *UserController) Login() {
	return
}

func (u UserController) GetUsers() {
	return
}

func (u UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		user, err := u.userRepo.GetUser(userId)

	}
}
