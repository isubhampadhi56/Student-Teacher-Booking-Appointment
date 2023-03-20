package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
)

func SearchTeacher(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("teacher")
	result, err := helper.SearchTeacherHelper(searchTerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&result)
}

func BookAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment model.Appointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "invalid request"})
	}
	appointment.Student = r.Context().Value("curuser").(string)
	email, name, err := helper.AppointmentHelper(&appointment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "unable to make appointments"})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ResponseStatus{Status: "appointment booked"})
	to := []string{email}
	subject := "Appointment Booked"
	body := fmt.Sprintf("Hi %v, %v has booked an appointment with you", name, appointment.Student)
	go helper.SendMail(to, subject, body)
}
