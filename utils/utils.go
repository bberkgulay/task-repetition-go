package utils

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// Headers set header to request
func Headers(r http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"})
	return handlers.CORS(headersOk, originsOk, methodsOk)(r)
}
