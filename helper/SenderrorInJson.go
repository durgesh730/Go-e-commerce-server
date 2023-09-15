package helper

import (
	"encoding/json"
	"net/http"
)

func SendJSONError(response http.ResponseWriter, message string, status int) {
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(status)
	json.NewEncoder(response).Encode(map[string]string{
		"error": message,
	})
}
