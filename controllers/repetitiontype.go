package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/bberkgulay/task-repetition-go/models"
	"github.com/bberkgulay/task-repetition-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c Controller) GetRepetitionTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		var repetitionTypes []models.RepetitionType

		cursor, err := c.DB.Collection("repetitiontypes").Find(context.TODO(), bson.D{})

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		if err = cursor.All(context.TODO(), &repetitionTypes); err != nil {
			error.Message = "Error while parsing data."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		utils.SendSuccess(w, repetitionTypes)
	}
}

func (c Controller) GetNextRepetitionType(repetitionTypeId primitive.ObjectID) (*models.RepetitionType, error) {
	var repetitionType *models.RepetitionType
	var nextRepetitionType models.RepetitionType

	order := 0

	if !repetitionTypeId.IsZero() {
		filter := bson.D{{Key: "_id", Value: repetitionTypeId}}

		err := c.DB.Collection("repetitiontypes").FindOne(context.TODO(), filter).Decode(&repetitionType)
		if err != nil {
			return nil, errors.New("error")
		}
		order = repetitionType.Order
	}

	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(bson.D{{Key: "order", Value: 1}})
	filter := bson.D{{Key: "order", Value: bson.M{"$gt": order}}}

	findErr := c.DB.Collection("repetitiontypes").FindOne(context.TODO(), filter, &queryOptions).Decode(&nextRepetitionType)
	if findErr != nil {
		return nil, nil
	}

	return &nextRepetitionType, nil
}
