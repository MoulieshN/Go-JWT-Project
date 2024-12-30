package helpers

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type signedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserType  string
	Uid       string
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, userType string, userId string) (string, string, error) {
	claims := &signedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Uid:       userId,
		UserType:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &signedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(168))),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return token, refresh_token, nil
}

func ValidateToken(signedToken string) (claims *signedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&signedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*signedDetails)
	if !ok {
		msg = "The token is invalid: " + err.Error()
		return
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		// Token is expired
		msg = "The token is expired: " + err.Error()
		return
	}

	return claims, msg
}
