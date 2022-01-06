package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bberkgulay/task-repetition-go/db"
	"github.com/bberkgulay/task-repetition-go/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database := db.Connect()
	fmt.Println(database)

	router := mux.NewRouter()

	srv := &http.Server{
		Handler:      utils.Headers(router), // Set header to routes
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Application is running at ", os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}
