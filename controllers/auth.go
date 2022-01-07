package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bberkgulay/task-repetition-go/models"
	"github.com/bberkgulay/task-repetition-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @route       POST /api/v1/auth/login
// @access      Public
func (c Controller) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var userOnDB models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" || user.Password == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		filter := bson.D{{Key: "email", Value: user.Email}}

		err := c.DB.Collection("users").FindOne(context.TODO(), filter).Decode(&userOnDB)
		if err != nil || !utils.CompareHashAndPassword(userOnDB.Password, user.Password) {
			error.Message = "Incorrect Email/Passwordasdasd"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, true)
	}
}

// @route       POST /api/v1/auth/register
// @access      Public
func (c Controller) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" || user.Name == "" || user.Surname == "" || user.Password == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		hashedPassword := utils.HashPassword(user.Password)
		if hashedPassword == "" {
			error.Message = "Error while hashing password."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		user.Password = hashedPassword

		filter := bson.D{{Key: "email", Value: user.Email}}
		existedUser, findError := c.DB.Collection("users").CountDocuments(context.TODO(), filter)

		if existedUser > 0 {
			error.Message = "That username is taken. Try another."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		if findError != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		insertResult, err := c.DB.Collection("users").InsertOne(context.TODO(), user)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, insertResult.InsertedID)
	}
}

// @description Middleware for authentication of endpoints.
func (c Controller) LoginControl(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/auth/") {
			h.ServeHTTP(w, r)
		} else {
			username, password, ok := r.BasicAuth()
			var error models.Error

			if !ok {
				error.Message = "No basic auth present"
				utils.SendError(w, http.StatusUnauthorized, error)
				return
			}

			if !isAuthorised(username, password, c.DB, r) {
				error.Message = "Invalid username or password"
				utils.SendError(w, http.StatusUnauthorized, error)
				return
			}

			h.ServeHTTP(w, r)
		}

	})
}

//@description Authorisation control of user if authorised, user id will be added to header.
func isAuthorised(username string, password string, db *mongo.Database, r *http.Request) bool {
	var user models.User

	filter := bson.D{{Key: "email", Value: username}}

	err := db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)

	if err != nil || !utils.CompareHashAndPassword(user.Password, password) {
		return false
	}

	r.Header.Set("userID", user.ID.Hex())

	return true

}
