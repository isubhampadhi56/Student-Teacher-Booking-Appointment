package controller

import (
	"encoding/json"
	"net/http"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
)

func bookAppointment(w http.ResponseWriter, r *http.Request) {

}

func SearchTeacher(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("teacher")
	result, err := helper.SearchTeacherHelper(searchTerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&result)
}
