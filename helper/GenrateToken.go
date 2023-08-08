package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT secret key to sign the tokens
var JwtScret = []byte("durgeshchaudhary")

func GererateToken(userId string) (string, error) {
	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // token will  expire in 24 hours
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(JwtScret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
