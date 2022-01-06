package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bberkgulay/task-repetition-go/controllers"
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
	controller := controllers.Controller{DB: database}

	router := mux.NewRouter()

	router.HandleFunc("/auth/register", controller.Register()).Methods("POST")
	router.HandleFunc("/auth/login", controller.Login()).Methods("POST")

	router.Use(controller.LoginControl)

	srv := &http.Server{
		Handler:      utils.Headers(router), // Set header to routes
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Application is running at ", os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}
