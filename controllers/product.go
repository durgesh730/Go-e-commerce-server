package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MaleCreateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var product models.MaleProduct
	_ = json.NewDecoder(r.Body).Decode(&product)

	result, _ := database.Productdata.InsertOne(context.Background(), product)
	product.Id = result.InsertedID.(primitive.ObjectID)

	json.NewEncoder(w).Encode(product.Id)
}

func FemaleCreateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var product models.FemaleProduct
	_ = json.NewDecoder(r.Body).Decode(&product)

	result, _ := database.Productdata.InsertOne(context.Background(), product)
	product.Id = result.InsertedID.(primitive.ObjectID)

	json.NewEncoder(w).Encode(product.Id)
}

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
