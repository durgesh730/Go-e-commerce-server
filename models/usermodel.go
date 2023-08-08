package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"Name"`
	Email    string `json:"email" bson:"Email"`
	Password string `json:"password" bson:"Password"`
}