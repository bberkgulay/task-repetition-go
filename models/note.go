package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Note      string             `bson:"note,omitempty"`
	Important bool               `bson:"important,omitempty"`
	User      primitive.ObjectID `bson:"user,omitempty"`
	Task      primitive.ObjectID `bson:"task,omitempty"`
}
