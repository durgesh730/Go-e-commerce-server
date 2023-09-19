package helper

import (
	"fmt"
	"net/http"

	"github.com/durgesh730/authenticationInGo/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetObjectIDFromToken(request *http.Request) (primitive.ObjectID, error) {
	val := request.Context().Value(middleware.UserIDKey)
	if val == nil {
		return primitive.NilObjectID, fmt.Errorf("No user ID present")
	}

	userId, ok := val.(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("User ID is of the wrong type")
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("Invalid user ID format")
	}

	return objectId, nil
}
