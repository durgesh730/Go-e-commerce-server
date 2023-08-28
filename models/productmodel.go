package models

import (
	// "time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaleProduct struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Brand     string             `bson:"brand"`
	Discount  int                `bson:"discount"`
	Size      []string           `bson:"size"`
	Fabric    string             `bson:"fabric"`
	Packof    int                `bson:"packof"`
	Fit       string             `bson:"fit"`
	Images    []string           `bson:"image"`
	Price     float64            `bson:"price"`
	Gender    string             `bson:"gender"`
	Type      string             `bson:"type"`
	StyleCode string             `bson:"styleCode"`
	NeckType  string             `bson:"NeckType"`
	Count     int                `bson:"count"`
	// CreatedAt time.Time          `bson:"created_at"`
	// UpdatedAt time.Time          `bson:"updated_at"`
}

type FemaleProduct struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Brand      string             `bson:"brand"`
	Discount   int                `bson:"discount"`
	SariLength string             `bson:"sariLength"`
	Fabric     string             `bson:"fabric"`
	Packof     int                `bson:"packof"`
	Weight     string             `bson:"weight"`
	Images     []string           `bson:"image"`
	Price      float64            `bson:"price"`
	Gender     string             `bson:"gender"`
	Type       string             `bson:"type"`
	StyleCode  string             `bson:"styleCode"`
	Pattern    string             `bson:"pattern"`
	SariStyle  string             `bson:"sariStyle"`
	Occasion   string             `bson:"occasion"`
	Count      int                `bson:"count"`

	// CreatedAt time.Time          `bson:"created_at"`
	// UpdatedAt time.Time          `bson:"updated_at"`
}
