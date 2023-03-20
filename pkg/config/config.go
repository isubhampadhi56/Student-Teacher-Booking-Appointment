package config

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setup mongoDB Connection

const Database = "student"

// const dbName = "studentDB"
// const colName = "users"

var Conn *mongo.Client
var EmailAuth smtp.Auth
var JwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
		return
	}
	JwtSecret = []byte(os.Getenv("JWT_KEY"))
	var connString = os.Getenv("CONNECTION_STRING")
	var email = os.Getenv("EMAIL")
	var emailPassword = os.Getenv("PASSWORD")
	clientOption := options.Client().ApplyURI(connString)
	conn, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		fmt.Println("Error connecting to Database")
		panic(err)
	}
	fmt.Println("Connection to MongoDB Successful")
	Conn = conn
	EmailAuth = smtp.PlainAuth("", email, emailPassword, "smtp.gmail.com")
}
