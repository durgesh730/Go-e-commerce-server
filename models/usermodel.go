package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	Street  string             `json:"street" bson:"street"`
	City    string             `json:"city" bson:"city"`
	State   string             `json:"state" bson:"state"`
	ZipCode string             `json:"zipCode" bson:"zipCode"`
	Country string             `json:"country" bson:"country"`
}

type User struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Addresses []Address          `json:"addresses" bson:"addresses"`
}
