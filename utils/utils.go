package utils

import (
	"encoding/json"
	"net/http"

	"github.com/bberkgulay/task-repetition-go/models"
	"github.com/gorilla/handlers"
	"golang.org/x/crypto/bcrypt"
)

// Headers set header to request
func Headers(r http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"})
	return handlers.CORS(headersOk, originsOk, methodsOk)(r)
}

func SendError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	// fmt.Println(data)
	json.NewEncoder(w).Encode(data)
}

func HashPassword(password string) string {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}

func CompareHashAndPassword(hashedPassword string, password string) bool {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)

	return err == nil
}
