package helper

import (
	"context"
	"log"

	"github.com/StudentTeacher-Booking-Appointment/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUnapprovedUsers() ([]bson.M, error) {
	const colName = "users"
	filter := bson.M{"isvalidated": false}
	projection := bson.M{"username": 1, "firstname": 1, "lastname": 1, "email": 1, "phone": 1, "class": 1, "stream": 1, "_id": 0}
	coll := config.Conn.Database(config.Database).Collection(colName)
	cursor, err := coll.Find(context.TODO(), filter, options.Find().SetProjection(projection))
	if err != nil {
		log.Println(err)
		return []bson.M{}, err
	}
	var results []bson.M
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Println(err)
			continue
		}
		results = append(results, result)
	}
	return results, nil
}

func ApproveHelper(user string) error {
	const colName = "users"
	filter := bson.M{"username": user, "isvalidated": false}
	update := bson.M{"$set": bson.M{"isvalidated": true}}
	coll := config.Conn.Database(config.Database).Collection(colName)
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
