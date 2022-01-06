package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RepetitionType struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Order int                `bson:"order,omitempty"`
	Day   int                `bson:"day,omitempty"`
}
