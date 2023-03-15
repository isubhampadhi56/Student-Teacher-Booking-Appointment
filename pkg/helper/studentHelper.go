package helper

import (
	"context"
	"log"

	"github.com/StudentTeacher-Booking-Appointment/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchTeacherHelper(searchTerm string) ([]bson.M, error) {
	pattern := primitive.Regex{Pattern: searchTerm, Options: "i"}
	filter := bson.M{
		"$or": []bson.M{
			{"firstname": bson.M{"$regex": pattern}, "role": "teacher"},
			{"lastname": bson.M{"$regex": pattern}, "role": "teacher"},
			{"dept": bson.M{"$regex": pattern}, "role": "teacher"},
			{"subject": bson.M{"$regex": pattern}, "role": "teacher"},
		},
	}
	projection := bson.M{"password": 0, "token": 0, "refresh_token": 0, "isvalidated": 0, "updatedat": 0, "role": 0, "phone": 0, "email": 0, "_id": 0}
	const colName = "users"
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
			continue
		}
		results = append(results, result)
	}
	return results, nil
}
