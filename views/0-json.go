package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error" example:"validation_error"`
	Message string `json:"message" example:"Phone number already registered"`
}

// SuccessResponse represents generic success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
	Status  string `json:"status" example:"success"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	RespondWithJSON(w, code, ErrorResponse{
		Error: fmt.Sprintf("%s: %v", msg, err),
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Message is set: %s", err)
	}
}
