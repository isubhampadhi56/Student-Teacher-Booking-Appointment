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
	"github.com/golang-jwt/jwt/v5"
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

var jwtSecret = []byte("Hello")

// Generate jwt Token
func GenerateToken(user map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user["username"].(string)
	claims["role"] = user["role"].(string)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["refresh_token"] = user["refresh_token"].(string)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//Validate jwt Token

func ValidateToken(tkn string) (bool, error) {
	if tkn == "" {
		return false, errors.New("missing authorization header.please login again")
	}
	tokenString := tkn
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return the key used to sign the token
		return jwtSecret, nil
	})
	if err != nil {
		return false, err
	}
	fmt.Println(token)
	return true, nil
}
