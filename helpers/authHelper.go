package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	return nil
}

func MatchUserTypeToUid(c *gin.Context, userID string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil
	if userType == "USER" && userID != uid {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return nil
}
