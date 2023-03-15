package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
)

// Write defination for all http Handler
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.Authenticate
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "invalid user data"})
		return
	}
	userFromDB, err := helper.LoginHelper(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: fmt.Sprint(err)})
		return
	}
	tokenString, err := helper.GenerateToken(userFromDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "unable to create token"})
	}
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 6),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    tokenString,
	})
	json.NewEncoder(w).Encode(model.ResponseStatus{Status: "Login Successful", Message: fmt.Sprintf("Hello %v. You are a %v", userFromDB["firstname"].(string), userFromDB["role"].(string))})

	// fmt.Println(validateToken(tokenString))
	// fmt.Printf("%+v", userDB)
}

func StudentSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student model.Students
	//var errorj model.ErrorObj
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "invalid user data"})
		return
	}
	err = helper.SignupHelper(&student)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: fmt.Sprint(err)})
		return
	} else {
		json.NewEncoder(w).Encode(model.ResponseStatus{Status: "user creation successful"})
	}

}
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    "",
	})
}
func Home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Authenticated")
}
