package helper

import (
	"log"
	"net/smtp"

	"github.com/StudentTeacher-Booking-Appointment/pkg/config"
)

func SendMail(to []string, subject string, body string) error {
	// Set up the message headers.
	from := "your.email@example.com"

	// Compose the message body.
	message := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send the email.
	err := smtp.SendMail("smtp.gmail.com:587", config.EmailAuth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
