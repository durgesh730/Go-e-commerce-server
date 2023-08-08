package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/router"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// connect database
	err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to mongoDB", err)
		return
	}
	// connect router
	r := router.Router()

	// Set up CORS
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "application/json"}, // You might need to add to this
	})
	handler := c.Handler(r)

	//start the server
	fmt.Println("server linstening on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
