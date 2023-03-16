package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/StudentTeacher-Booking-Appointment/pkg/helper"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
	"github.com/go-chi/jwtauth/v5"
)

func RequireAuth(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		tkn := jwtauth.TokenFromCookie(r)
		if tkn == "" {
			log.Println("user not loggedin")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(model.ResponseStatus{Err: "user not loggedin"})
			return
		}
		// fmt.Println(tkn)
		// ok,err :=
		token, err := helper.ValidateToken(tkn)
		if err != nil {
			// 	// If not, return a 401 Unauthorized error
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(model.ResponseStatus{Err: "session expired please relogin"})
			return
		}
		if role == token["role"] {
			fmt.Println(r.URL.User.Username())
			next(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(model.ResponseStatus{Err: "insufficient permision to access this resource"})
		// If the user is authenticated, call the next handler function

	}
}
