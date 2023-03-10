package main

import (
	"fmt"
	"net/http"

	"github.com/StudentTeacher-Booking-Appointment/pkg/routes"
	"github.com/go-chi/chi/v5"
)

var PORT = "3000"

func main() {
	var Router = chi.NewRouter()
	routes.RegisterRoute(Router)
	fmt.Printf("Starting server on Port %v", PORT)
	err := http.ListenAndServe(":"+PORT, Router)
	if err != nil {
		panic(err)
	}
}
