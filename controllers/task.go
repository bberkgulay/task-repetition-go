package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bberkgulay/task-repetition-go/models"
	"github.com/bberkgulay/task-repetition-go/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @route       POST /api/v1/tasks
// @access      Private
// @description Adds task for user.
func (c Controller) AddTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		var error models.Error

		json.NewDecoder(r.Body).Decode(&task)

		if task.Title == "" {
			error.Message = "Enter missing fields. (Title)"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		//Getting user from header.
		userId, hexError := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if hexError != nil {
			error.Message = "Error occurred getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		task.User = userId

		insertResult, err := c.DB.Collection("tasks").InsertOne(context.TODO(), task)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, insertResult.InsertedID)
	}
}

// @route       GET /api/v1/tasks
// @access      Private
// @description Returns tasks of user.
func (c Controller) GetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		var tasks []models.Task

		//Getting user from header.
		userId, hexError := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if hexError != nil {
			error.Message = "Error occurred about user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "user", Value: userId}}

		cursor, err := c.DB.Collection("tasks").Find(context.TODO(), filter)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		if err = cursor.All(context.TODO(), &tasks); err != nil {
			error.Message = "Error while parsing data."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, tasks)
	}
}

// @route       GET /api/v1/tasks/{id}
// @access      Private
// @description Gets task by id.
func (c Controller) GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var task models.Task
		var error models.Error

		params := mux.Vars(r)

		objectId, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			error.Message = "Incorrect ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		//Getting user from header for owner control.
		userId, err := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if err != nil {
			error.Message = "Error while getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "_id", Value: objectId}, {Key: "user", Value: userId}}
		findError := c.DB.Collection("tasks").FindOne(context.TODO(), filter).Decode(&task)

		if findError != nil {
			error.Message = "Server Error."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, task)
	}
}

// @route       PUT /api/v1/tasks/{id}
// @access      Private
// @description Updates task by id with owner control.
func (c Controller) UpdateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		var error models.Error

		json.NewDecoder(r.Body).Decode(&task)

		if task.Title == "" {
			error.Message = "Enter missing fields.(Title)"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		params := mux.Vars(r)

		objectId, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			error.Message = "Incorrect ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		//Getting user from header for owner control.
		userId, err := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if err != nil {
			error.Message = "Error while getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "_id", Value: objectId}, {Key: "user", Value: userId}}

		result, err := c.DB.Collection("tasks").UpdateOne(
			context.TODO(),
			filter,
			bson.D{
				{"$set", task},
			},
		)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, result)
	}
}

// @route       DELETE /api/v1/tasks/{id}
// @access      Private
// @description Deletes task by id with owner control.
func (c Controller) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)

		objectId, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			error.Message = "Incorrect ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		//Getting user from header for owner control.
		userId, err := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if err != nil {
			error.Message = "Error while getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "_id", Value: objectId}, {Key: "user", Value: userId}}

		result := c.DB.Collection("tasks").FindOneAndDelete(context.TODO(), filter).Err()

		if result != nil {
			error.Message = "No task to delete"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		utils.SendSuccess(w, "Successful")
	}
}

// @route       PUT /api/v1/tasks/{id}/complete
// @access      Private
// @description Completes task and finds suitable repetition type and repetition date.
func (c Controller) CompleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var task models.Task
		var error models.Error

		params := mux.Vars(r)

		objectId, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			error.Message = "Incorrect ID value."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		userId, err := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if err != nil {
			error.Message = "Error while getting user."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{"_id", objectId}, {Key: "user", Value: userId}}
		findError := c.DB.Collection("tasks").FindOne(context.TODO(), filter).Decode(&task)

		if findError != nil {
			error.Message = "Server Error."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		if !task.CompletedDay.IsZero() {
			error.Message = "Task is already completed."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		//TODO begin day control
		// if task.RepetitionBeginDay > time.Now() {
		// 	error.Message = ""
		// 	utils.SendError(w, http.StatusBadRequest, error)
		// 	return
		// }

		nextRepetitionType, err := c.GetNextRepetitionType(task.RepetitionType)
		if err != nil {
			error.Message = "Error"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		message := "Successful"

		//no next repetition day means the user has completed the task.
		if nextRepetitionType == nil {
			task.CompletedDay = time.Now()
			message = "Task is completed successfully."
		} else {
			//Sets next repetition begin date and type.
			task.RepetitionBeginDay = time.Now().AddDate(0, 0, nextRepetitionType.Day)
			task.RepetitionType = nextRepetitionType.ID
			message = "Successful"
		}

		filter = bson.D{{Key: "_id", Value: objectId}, {Key: "user", Value: userId}}
		result, err := c.DB.Collection("tasks").UpdateOne(
			context.TODO(),
			filter,
			bson.D{
				{"$set", task},
			})
		_ = result

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, message)
	}
}
