package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setup mongoDB Connection
const user = "subham"
const password = "hello123"
const Database = "student"
const connString = "mongodb+srv://" + user + ":" + password + "@dev-01.lofuxvh.mongodb.net/?retryWrites=true&w=majority"

// const dbName = "studentDB"
// const colName = "users"

var Conn *mongo.Client

func init() {
	clientOption := options.Client().ApplyURI(connString)
	conn, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		fmt.Println("Error connecting to Database")
		panic(err)
	}
	fmt.Println("Connection to MongoDB Successful")
	Conn = conn
}
