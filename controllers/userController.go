package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MoulieshN/Go-JWT-Project.git/helpers"
	"github.com/MoulieshN/Go-JWT-Project.git/models"
	"github.com/MoulieshN/Go-JWT-Project.git/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) UserController {
	return UserController{userRepo: repo}
}

func HashPassword(userPassword string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(hashedPassword)
}

func VerfiyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Email or password is incorrect"
		check = false
	}
	return check, msg
}

func (u *UserController) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Converting password into hashed password for more security
		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

		// Save user in a db
		incrementNo, err := u.userRepo.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate the token and refresh token
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserId)

		// update the user with token and refresh token
		err = u.userRepo.UpdateUserToken(token, refreshToken, user.UserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": incrementNo})
	}
}

func (u *UserController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the user by email
		foundUser, err := u.userRepo.GetUserByEmail(*user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		// check the login user's password and saved user password is save
		isVerified, msg := VerfiyPassword(*user.Password, *foundUser.Password)
		if !isVerified {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)

		// update the user with token and refresh token
		err = u.userRepo.UpdateUserToken(token, refreshToken, user.UserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": foundUser})
	}
}

func (u UserController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"Error": err.Error()})
			return
		}

		itemsPerPage, err := strconv.Atoi(c.Query("itemsPerPage"))
		if err != nil || itemsPerPage < 1 {
			itemsPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		offset := (page - 1) * itemsPerPage

		users, err := u.userRepo.GetUsers(page, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data := make(map[string]interface{})
		data["data"] = users

		c.JSON(http.StatusOK, data)
	}
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
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}
