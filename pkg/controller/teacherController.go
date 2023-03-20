package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListAppointments(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("curuser").(string)
	appointments, err := helper.GetAllAppointments(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "internal server error"})
	}
	json.NewEncoder(w).Encode(&appointments)
}

func ApproveAppoitnments(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("curuser").(string)
	var postData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&postData)
	id, _ := primitive.ObjectIDFromHex(postData["id"].(string))
	err := helper.ApproveAppointmentHelper(id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: fmt.Sprint(err)})
		return
	}
	json.NewEncoder(w).Encode(model.ResponseStatus{Status: "appointment approved"})
}

func ScheduleAppointment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("curuser").(string)
	var postData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&postData)
	id, _ := primitive.ObjectIDFromHex(postData["id"].(string))
	time, _ := time.Parse("2006-01-02 15:04", postData["shedtime"].(string))
	appointmentDetails, err := helper.ScheduleAppointmentHelper(id, user, time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: fmt.Sprint(err)})
		return
	}
	json.NewEncoder(w).Encode(appointmentDetails)
	userDetails, err := helper.FindUser(appointmentDetails["student"].(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "unable to send email.But meeting has been scheduled."})
		return
	}
	body := fmt.Sprintf("Hi %v, Your Appointment has been scheduled with %v on %v \nPlease Join the meeting with the link %v ", appointmentDetails["student"].(string), appointmentDetails["teacher"].(string), appointmentDetails["shedtime"].(primitive.DateTime).Time().String(), appointmentDetails["meetlink"].(string))
	to := []string{userDetails.Email}
	subject := "Appointment Scheduled"
	go helper.SendMail(to, subject, body)
}
