package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string
	Price      float64
	Videocount int
}
