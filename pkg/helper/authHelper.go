package helper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/StudentTeacher-Booking-Appointment/pkg/config"
	"github.com/StudentTeacher-Booking-Appointment/pkg/model"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Responsible from authenticate user from DB
func LoginHelper(user model.Authenticate) (map[string]interface{}, error) {
	// user.Username
	const colName = "users"
	coll := config.Conn.Database(config.Database).Collection(colName)
	filter := bson.M{"username": user.Username}
	userFromDB := map[string]interface{}{}
	err := coll.FindOne(context.TODO(), filter).Decode(&userFromDB)
	// defer cursor.Close(context.TODO())
	if err == mongo.ErrNoDocuments {
		log.Println("username not exist")
		return userFromDB, errors.New("username does not exist")
	}
	//DBUser := userFromDB.(Users)

	fmt.Printf("%+v", userFromDB)
	inputPassword, ok := userFromDB["password"].(string)
	if !ok {
		fmt.Println("Error: unable to extract password value")
	}
	err = bcrypt.CompareHashAndPassword([]byte(inputPassword), []byte(user.Password))
	if err == nil {
		return userFromDB, nil
	}
	log.Println(err)
	return userFromDB, errors.New("username or password incorrect")
}

//Signup Helper responsible for adding basic Attributes to objects and write it to DB

func SignupHelper(user model.User) error {
	user.SetRole()
	err := user.Validate()
	if err != nil {
		return err
	}
	err = user.EncryptPassword()
	if err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}
	fmt.Println(user)
	const colName = "users"
	coll := config.Conn.Database(config.Database).Collection(colName)

	//Creating Index username,phone,email unique
	func() {
		index := mongo.IndexModel{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		}
		_, err = coll.Indexes().CreateOne(context.TODO(), index)
		if err != nil {
			log.Println(err)
		}
		index = mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		}
		_, err = coll.Indexes().CreateOne(context.TODO(), index)
		if err != nil {
			log.Println(err)
		}
		index = mongo.IndexModel{
			Keys:    bson.M{"phone": 1},
			Options: options.Index().SetUnique(true),
		}
		_, err = coll.Indexes().CreateOne(context.TODO(), index)
		if err != nil {
			log.Println(err)
		}
	}()
	//Inserting into Database
	id, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		if merr, ok := err.(mongo.WriteException); ok {
			for _, we := range merr.WriteErrors {
				switch we.Code {
				case 11000:
					if strings.Contains(we.Message, "username") {
						return errors.New("username already exist")
					} else if strings.Contains(we.Message, "email") {
						return errors.New("email already exist")
					} else if strings.Contains(we.Message, "phone") {
						return errors.New("phone already exist")
					}
				}
			}
		}
		log.Println(err)
		return errors.New("error creating user")
	}
	log.Printf("user created %v", id.InsertedID)
	return nil
}

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", config.JwtSecret, nil)
}

// Generate jwt Token
func GenerateToken(user map[string]interface{}) (string, error) {

	claims := map[string]interface{}{
		"username":      user["username"].(string),
		"role":          user["role"].(string),
		"exp":           time.Now().Add(time.Minute * 10).Unix(),
		"refresh_token": user["refresh_token"].(string),
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", errors.New("internal server error")
	}
	return tokenString, nil
}

//Validate jwt Token

func ValidateToken(tokenString string) (map[string]interface{}, error) {

	token, err := jwtauth.VerifyToken(tokenAuth, tokenString)
	if err != nil {
		return map[string]interface{}{}, err
	}
	tokenMap, err := token.AsMap(context.TODO())
	if err != nil {
		return map[string]interface{}{}, errors.New("internal server error")
	}
	return tokenMap, nil
}

func FindUser(user string) (model.Users, error) {
	const colName = "users"
	coll := config.Conn.Database(config.Database).Collection(colName)
	var userDetails model.Users
	filter := bson.M{"username": user}
	err := coll.FindOne(context.TODO(), filter).Decode(&userDetails)
	userDetails.RemovePassword()
	if err != nil {
		return model.Users{}, err
	}
	return userDetails, nil
}
