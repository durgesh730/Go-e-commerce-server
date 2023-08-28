package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ProductId primitive.ObjectID `bson:"productId"`
	UserId    primitive.ObjectID `bson:"userId"`
	TotalItem int                `bson:"totalItem"`
}
