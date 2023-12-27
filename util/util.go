package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiError struct {
	Field string
	Msg   string
}

type ApiResponse struct {
	Data   interface{} `json:"data"`
	Errors []ApiError  `json:"errors"`
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}, errors []ApiError) {

	data := ApiResponse{Data: payload,Errors: errors}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Printf("Failed to marshal %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
}

