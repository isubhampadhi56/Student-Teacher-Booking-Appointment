package helper

import (
	"context"
	"log"

	"github.com/StudentTeacher-Booking-Appointment/pkg/config"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
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

func AppointmentHelper(appointment *model.Appointment) (string, string, error) {
	var collName = "users"
	coll := config.Conn.Database(config.Database).Collection(collName)
	filter := bson.M{"username": appointment.Teacher, "role": "teacher"}
	var user model.Users
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Panicln(err)
		return "", "", err
	}
	collName = "appointment"
	coll = config.Conn.Database(config.Database).Collection(collName)
	appointment.Book()
	_, err = coll.InsertOne(context.TODO(), appointment)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	return user.Email, user.FirstName, nil
}
