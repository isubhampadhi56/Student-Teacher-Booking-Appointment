package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	Id            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Teacher       string             `json:"teacher" bson:"teacher" validate:"required"`
	Student       string             `json:"student" bson:"student" validate:"required"`
	PreferedTime  time.Time          `json:"preftime" bson:"preftime"`
	Isconfirmed   bool               `json:"isconfirmed" bson:"isconfirmed"`
	ScheduledTime time.Time          `json:"shedtime" bson:"shedtime"`
	Status        string             `json:"status" bson:"status"`
	MeetingLink   string             `json:"meetlink" bson:"meetlink"`
	CreatedAt     time.Time          `json:"createdat" bson:"createdat"`
}

func (app *Appointment) Book() {
	app.Id = primitive.NewObjectID()
	app.Isconfirmed = false
	app.Status = "requested"
	app.MeetingLink = ""
	app.CreatedAt = time.Now()
}
