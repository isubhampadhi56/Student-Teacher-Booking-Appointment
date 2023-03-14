package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		tkn, err := r.Cookie("token")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ResponseStatus{Err: "user not loggedin"})
			return
		}
		ok, err := helper.ValidateToken(tkn.Value)
		if !ok {
			// If not, return a 401 Unauthorized error
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ResponseStatus{Err: "session failed please relogin"})
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		// If the user is authenticated, call the next handler function
		next(w, r)
	}
}
