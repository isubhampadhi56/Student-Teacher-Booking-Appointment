package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
)

func AddTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher model.Teachers
	err := json.NewDecoder(r.Body).Decode(&teacher)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "invalid request body"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = helper.SignupHelper(&teacher)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(model.ResponseStatus{Status: "teacher added"})
	w.WriteHeader(http.StatusOK)
}

func ListUnApprovedStudent(w http.ResponseWriter, r *http.Request) {
	result, err := helper.GetUnapprovedUsers()
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(&result)
}

func ApproveStudent(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "invalid request"})
		w.WriteHeader(http.StatusBadRequest)
	}
	username := user["username"].(string)
	err = helper.ApproveHelper(username)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "unable to approve user"})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
