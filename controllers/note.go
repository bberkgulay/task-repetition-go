package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bberkgulay/task-repetition-go/models"
	"github.com/bberkgulay/task-repetition-go/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c Controller) AddNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note models.Note
		var error models.Error

		json.NewDecoder(r.Body).Decode(&note)

		if note.Note == "" {
			error.Message = "Enter missing fields. (Note)"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		userId, hexError := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if hexError != nil {
			error.Message = "Error occured getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		note.User = userId

		params := mux.Vars(r)

		taskId, err := primitive.ObjectIDFromHex(params["task_id"])
		if err != nil {
			error.Message = "Incorrect Task ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		note.Task = taskId

		//task kullanıcının mı?
		filter := bson.D{{Key: "user", Value: userId}, {Key: "_id", Value: taskId}}
		taskCountOfUser, findError := c.DB.Collection("tasks").CountDocuments(context.TODO(), filter)

		if taskCountOfUser == 0 {
			error.Message = "Error"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		if findError != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		insertResult, err := c.DB.Collection("notes").InsertOne(context.TODO(), note)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, insertResult.InsertedID)
	}
}

func (c Controller) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		var notes []models.Note

		userId, hexError := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if hexError != nil {
			error.Message = "Error occured about notes."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		params := mux.Vars(r)

		taskId, err := primitive.ObjectIDFromHex(params["task_id"])
		if err != nil {
			error.Message = "Incorrect Task ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "user", Value: userId}, {Key: "task", Value: taskId}}

		cursor, err := c.DB.Collection("notes").Find(context.TODO(), filter)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		if err = cursor.All(context.TODO(), &notes); err != nil {
			error.Message = "Error while parsing data."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, notes)
	}
}

func (c Controller) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)

		userId, err := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if err != nil {
			error.Message = "Error while getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			error.Message = "Incorrect ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "user", Value: userId}, {Key: "_id", Value: id}}

		result := c.DB.Collection("notes").FindOneAndDelete(context.TODO(), filter).Err()

		if result != nil {
			error.Message = "No document to delete"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		utils.SendSuccess(w, "Successful")
	}
}
