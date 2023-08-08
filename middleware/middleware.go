package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// JWT secret key to sign the tokens
var jwtKey = []byte("durgeshchaudhary")
type Claims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}

type ContextKey string
const UserIDKey ContextKey = "userId"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			fmt.Println("Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

}
