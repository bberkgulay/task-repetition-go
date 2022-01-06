package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title,omitempty"`
	Link    string             `bson:"link,omitempty"`
	Summary string             `bson:"summary,omitempty"`
	Tags    []string           `bson:"tags,omitempty"`
	User    primitive.ObjectID `bson:"user,omitempty"`
}
