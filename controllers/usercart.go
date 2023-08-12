package controllers

import (
	"fmt"
	"net/http"

	"github.com/durgesh730/authenticationInGo/middleware"
)

func CreateProductCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	val := request.Context().Value(middleware.UserIDKey)
	if val == nil {
		// Handle the case where no userId is present
		http.Error(response, "No user ID present", http.StatusUnauthorized)
		return
	}

	userId, ok := val.(string)
	if !ok {
		// Handle the case where the userId is not a string
		http.Error(response, "User ID is of the wrong type", http.StatusInternalServerError)
		return
	}

	fmt.Println(userId)

}

func GetProductFromCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
}
func UpdateProductFromCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
}

func DeleteProductFromCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
}
