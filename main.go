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

	// endpoints are defined and mapped with functions.
	router.HandleFunc("/auth/register", controller.Register()).Methods("POST")
	router.HandleFunc("/auth/login", controller.Login()).Methods("POST")

	router.HandleFunc("/tasks", controller.GetTasks()).Methods("GET")
	router.HandleFunc("/tasks", controller.AddTask()).Methods("POST")
	router.HandleFunc("/tasks/{id}", controller.GetTask()).Methods("GET")
	router.HandleFunc("/tasks/{id}", controller.UpdateTask()).Methods("PUT")
	router.HandleFunc("/tasks/{id}", controller.DeleteTask()).Methods("DELETE")

	router.HandleFunc("/tasks/{id}/complete", controller.CompleteTask()).Methods("PUT")

	router.HandleFunc("/tasks/{task_id}/notes", controller.AddNote()).Methods("POST")
	router.HandleFunc("/tasks/{task_id}/notes", controller.GetNotes()).Methods("GET")
	router.HandleFunc("/notes/{id}", controller.DeleteNote()).Methods("DELETE")

	router.HandleFunc("/repetitiontypes", controller.GetRepetitionTypes()).Methods("GET")

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
