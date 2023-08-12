package controllers

import (
	"context"
	"encoding/json"

	// "fmt"
	"net/http"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// create products for male
func MaleCreateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var product models.MaleProduct
	_ = json.NewDecoder(r.Body).Decode(&product)

	result, _ := database.Productdata.InsertOne(context.Background(), product)
	product.Id = result.InsertedID.(primitive.ObjectID)

	json.NewEncoder(w).Encode(product.Id)
}

// create products for female
func FemaleCreateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var product models.FemaleProduct
	_ = json.NewDecoder(r.Body).Decode(&product)

	result, _ := database.Productdata.InsertOne(context.Background(), product)
	product.Id = result.InsertedID.(primitive.ObjectID)

	json.NewEncoder(w).Encode(product.Id)
}

// Fetch all Products
func GetProducts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var products []models.MaleProduct
	data, err := database.Productdata.Find(context.Background(), bson.M{})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	for data.Next(context.Background()) {
		var product models.MaleProduct
		data.Decode(&product)
		products = append(products, product)
	}

	if err := data.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(products)
}

// Fetch Products According to male and female according to query
func GetQueryProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	filter := bson.M{"gender": primitive.Regex{Pattern: query, Options: "i"}}
	curr, err := database.Productdata.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var results []models.MaleProduct
	for curr.Next(context.Background()) {
		var item models.MaleProduct
		curr.Decode(&item)
		results = append(results, item)
	}
	json.NewEncoder(w).Encode(results)
}

// get products by id

func GetProductsbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		http.Error(w, "Invalid query parameter 'q'", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": id}

	curr, err := database.Productdata.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return

	}
	var results []models.MaleProduct
	for curr.Next(context.Background()) {
		var item models.MaleProduct
		curr.Decode(&item)
		results = append(results, item)
	}

	json.NewEncoder(w).Encode(results)
}
