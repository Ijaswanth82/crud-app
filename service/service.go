package service

import (
	"context"
	"databaseconnection/model"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb+srv://eshwarpendem:eshwarpendem@cluster0.u7cmodf.mongodb.net/?retryWrites=true&w=majority"
var dbName string = "smurfDB"
var collectionName string = "courses"
var collection *mongo.Collection

func init() {
	fmt.Println("Inside init")
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")
	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("Creation of database and table success")
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping database success!!!")
}
func FindAll() ([]model.Course, error) {
	var courses []model.Course
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var course model.Course
		if err = cur.Decode(&course); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}
func AddOneRecord(course model.Course) error {
	course.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), course)
	if err != nil {
		return err
	}
	return nil
}
func DeleteRecord(courseName string, courses *[]model.Course) (int64, error) {
	cur, err := collection.DeleteMany(context.TODO(), bson.D{{Key: "name", Value: courseName}})
	if err != nil {
		return 0, err
	}
	return cur.DeletedCount, err
}
func FindByName(courseName string, courses *[]model.Course) error {
	cur, err := collection.Find(context.TODO(), bson.D{{Key: "name", Value: courseName}})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var course model.Course
		if err = cur.Decode(&course); err != nil {
			return err
		}
		(*courses) = append((*courses), course)
	}
	return nil
}
func UpdateByName(courseName string, course model.Course) (int64, error) {
	filter := bson.D{{Key: "name", Value: courseName}}
	if course.Name == "" && course.Price == 0 && course.Videocount == 0 {
		return 0, fmt.Errorf("course request body %+v is nil", course)
	}
	updatevalues := bson.D{}
	if course.Name != "" {
		updatevalues = append(updatevalues, primitive.E{Key: "name", Value: course.Name})
	}
	if course.Price != 0 {
		updatevalues = append(updatevalues, primitive.E{Key: "price", Value: course.Price})
	}
	if course.Videocount != 0 {
		updatevalues = append(updatevalues, primitive.E{Key: "videocount", Value: course.Videocount})
	}
	update := bson.D{{Key: "$set", Value: updatevalues}}
	res, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
