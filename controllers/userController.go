package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/helper"
	"github.com/durgesh730/authenticationInGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Parse the user registration data from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	// Check if the user already exists in MongoDB
	filter := bson.M{"email": user.Email}
	count, err := database.SaveData.CountDocuments(context.Background(), filter)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if count > 0 {
		// Respond with the success message
		response := make(map[string]interface{})
		response["msg"] = "User already exist"
		response["status"] = http.StatusConflict
		json.NewEncoder(w).Encode(response)
		return
	} else {
		// Hash the user's password before storing it
		hashedPassword, err := helper.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Failed to insert data into MongoDB", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
         
		if user.Addresses == nil {
			user.Addresses = []models.Address{}
		}		
		// Insert the user data into MongoDB
		Userdata, _ := database.SaveData.InsertOne(context.Background(), user)
		user.Id = Userdata.InsertedID.(primitive.ObjectID)

		// fmt.Fprintf(w, "The user ID is: %s", Userdata.InsertedID)
		tokenString, _ := helper.GererateToken(user.Id.Hex())

		//cookie
		cookie := http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 24), // token will expire in 24 hours
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		// Respond with the success message
		response := make(map[string]interface{})
		response["token"] = tokenString
		response["_id"] = Userdata
		response["msg"] = "user successfully register"
		response["status"] = http.StatusCreated

		json.NewEncoder(w).Encode(response)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Parse the user registration data from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	// find the user by email  in mongodb
	var exist models.User
	filter := bson.M{"email": user.Email}
	err = database.SaveData.FindOne(context.Background(), filter).Decode(&exist)
	if err != nil {
		response := make(map[string]interface{})
		response["msg"] = "Email not found"
		response["status"] = http.StatusNotFound
		json.NewEncoder(w).Encode(response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(user.Password))
	if err != nil {
		response := make(map[string]interface{})
		response["msg"] = "Invalid password"
		response["status"] = http.StatusUnauthorized
		json.NewEncoder(w).Encode(response)
		return
	}

	// generate token
	tokenString, _ := helper.GererateToken(exist.Id.Hex())
	//cookie
	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // token will expire in 24 hours
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// Respond with the success message
	response := make(map[string]interface{})
	response["token"] = tokenString
	response["msg"] = "User successfully logged In"
	response["user"] = exist
	response["status"] = http.StatusCreated

	json.NewEncoder(w).Encode(response)
}

func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	objectId, err := helper.GetObjectIDFromToken(request)
	if err != nil {
		switch err.Error() {
		case "No user ID present":
			helper.SendJSONError(response, err.Error(), http.StatusUnauthorized)
		case "User ID is of the wrong type", "Invalid user ID format":
			helper.SendJSONError(response, err.Error(), http.StatusInternalServerError)
		default:
			helper.SendJSONError(response, "Unexpected error", http.StatusInternalServerError)
		}
		return
	}

	filter := bson.M{"_id": objectId}
	var user models.User
	err = database.SaveData.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(response, `{"message": "User not found"}`, http.StatusNotFound)
		} else {
			http.Error(response, `{"message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(response).Encode(user); err != nil {
		http.Error(response, `{"message": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}


func AddAddress(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	objectId, err := helper.GetObjectIDFromToken(request)
	if err != nil {
		switch err.Error() {
		case "No user ID present":
			helper.SendJSONError(response, err.Error(), http.StatusUnauthorized)
		case "User ID is of the wrong type", "Invalid user ID format":
			helper.SendJSONError(response, err.Error(), http.StatusInternalServerError)
		default:
			helper.SendJSONError(response, "Unexpected error", http.StatusInternalServerError)
		}
		return
	}

	var address models.Address
	err = json.NewDecoder(request.Body).Decode(&address)
	if err != nil {
		helper.SendJSONError(response, "Address Not found", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$push": bson.M{"addresses": address}}
	
    fmt.Println(update)

	save, SaveErr := database.SaveData.UpdateOne(context.Background(), filter, update)
	if SaveErr != nil {
		http.Error(response, SaveErr.Error(), http.StatusInternalServerError)
		return
	}

	if save.MatchedCount == 0 {
		helper.SendJSONError(response, "User not found", http.StatusNotFound)
		return
	}

	responseJSON := map[string]interface{}{
		"message": "Document updated successfully",
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(responseJSON)
}
