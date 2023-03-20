package helper

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/StudentTeacher-Booking-Appointment/pkg/config"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllAppointments(user string) ([]model.Appointment, error) {
	const collName = "appointment"
	coll := config.Conn.Database(config.Database).Collection(collName)
	filter := bson.M{"teacher": user}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return []model.Appointment{}, errors.New("internal server error")
	}
	var appointments []model.Appointment
	for cursor.Next(context.TODO()) {
		var appointment model.Appointment
		cursor.Decode(&appointment)
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}
func ApproveAppointmentHelper(id primitive.ObjectID, user string) error {
	const collName = "appointment"
	coll := config.Conn.Database(config.Database).Collection(collName)
	filter := bson.M{"_id": id, "teacher": user}
	update := bson.M{"$set": bson.M{"isconfirmed": true}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}

func ScheduleAppointmentHelper(id primitive.ObjectID, user string, shedtime time.Time) (map[string]interface{}, error) {
	const collName = "appointment"
	coll := config.Conn.Database(config.Database).Collection(collName)
	filter := bson.M{"_id": id, "teacher": user}
	update := bson.M{"$set": bson.M{"shedtime": shedtime, "status": "scheduled"}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{}, errors.New("internal server error")
	}
	var result map[string]interface{}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
		return bson.M{}, errors.New("internal server error")
	}
	return result, nil
}
