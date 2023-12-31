package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/middleware"
	"github.com/durgesh730/authenticationInGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProductCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	val := r.Context().Value(middleware.UserIDKey)
	if val == nil {
		// Handle the case where no userId is present
		http.Error(w, "No user ID present", http.StatusUnauthorized)
		return
	}

	userId, ok := val.(string)
	if !ok {
		// Handle the case where the userId is not a string
		http.Error(w, "User ID is of the wrong type", http.StatusInternalServerError)
		return
	}
	var cart models.Cart
	_ = json.NewDecoder(r.Body).Decode(&cart)

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Token Not Valid", http.StatusBadRequest)
	}

	cart.UserId = id
	save, _ := database.Cartdata.InsertOne(context.Background(), cart)

	fmt.Println(save, "id")
	json.NewEncoder(w).Encode(save)
}

func GetProductFromCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	val := request.Context().Value(middleware.UserIDKey)
	if val == nil {
		// handle the case when user id not found
		http.Error(response, "No user ID present", http.StatusUnauthorized)
		return
	}

	userId, ok := val.(string)
	if !ok {
		http.Error(response, "User ID is worng type", http.StatusInternalServerError)
		return
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(response, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	filter := bson.M{"userId": objectId}
	curr, err := database.Cartdata.Find(request.Context(), filter)

	if err != nil {
		http.Error(response, "Internal server error", http.StatusInternalServerError)
		return
	}

	var cart []models.Cart
	for curr.Next(context.Background()) {
		var item models.Cart
		curr.Decode(&item)
		cart = append(cart, item)
	}

	if err := curr.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		http.Error(response, "Database iteration error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(response).Encode(cart); err != nil {
		http.Error(response, "Failed to encode response", http.StatusInternalServerError)
	}
}

func DeleteProductFromCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	query := request.URL.Query().Get("q")
	if query == " " {
		http.Error(response, "Missing query parameter 'q' ", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		http.Error(response, "Invalid query parameter 'q'", http.StatusBadRequest)
		return
	}

	filter := bson.M{"productId": id}
	curr := database.Cartdata.FindOneAndDelete(request.Context(), filter)

	if curr.Err() != nil {
		http.Error(response, curr.Err().Error(), http.StatusInternalServerError)
		return
	}

	var deletedProduct models.Cart
	if err := curr.Decode(&deletedProduct); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(deletedProduct)
}

func UpdateProductFromCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	query := request.URL.Query().Get("q")
	if query == " " {
		http.Error(response, "Invalid query paramter 'q'", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		http.Error(response, "Invalid Conversions of query", http.StatusBadRequest)
		return
	}

	// construct struct because of count is common in famle and male models
	var totalItem models.Cart
	_ = json.NewDecoder(request.Body).Decode(&totalItem)

	filter := bson.M{"productId": id}
	update := bson.M{"$set": bson.M{"totalItem": totalItem.TotalItem}}

	fmt.Println( totalItem.TotalItem, "filter and update")

	result, updateErr := database.Cartdata.UpdateOne(context.Background(), filter, update)
	if updateErr != nil {
		http.Error(response, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(response, "Product not found", http.StatusNotFound)
		return
	}

	// Send a success response
	responseJSON := map[string]interface{}{
		"message": "Document updated successfully",
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(responseJSON)
}


