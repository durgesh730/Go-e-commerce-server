package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/helper"
	"github.com/durgesh730/authenticationInGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/bson/primitive"
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
	filter := bson.M{"Email": user.Email}
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
		user.Password = hashedPassword;

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
	filter := bson.M{"Email": user.Email}
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
	response.Header().Set("content-type", "application/json")
	var users []models.User
	cursor, err := database.SaveData.Find(context.Background(), bson.M{}) // bson.M{} will match all documents in the collection
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for cursor.Next(context.Background()) {
		var user models.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(users)
}