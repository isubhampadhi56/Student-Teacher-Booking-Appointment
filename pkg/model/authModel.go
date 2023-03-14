package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Define Datamodels
type Users struct {
	Uid          primitive.ObjectID `json:"_" bson:"_id,omitempty"`
	Username     string             `json:"user" bson:"username,omitempty" validate:"required,min=4,max=50"`
	Password     string             `json:"password" bson:"password,omitempty" validate:"required,min=6,max=50"`
	Role         string             `json:"role" bson:"role,omitempty" validate:"required,eq=student|eq=teacher|eq=admin"`
	FirstName    string             `json:"firstname,omitempty" bson:"firstname,omitempty" validate:"required,min=2,max=100"`
	LastName     string             `json:"lastname,omitempty" bson:"lastname,omitempty" validate:"required,min=2,max=100"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	Phone        string             `json:"phone,omitempty" bson:"phone,omitempty" validate:"required"`
	IsValidated  bool               `json:"isvalidated" bson:"isvalidated"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAT    time.Time          `json:"updated_at"`
}

type Students struct {
	Users  `bson:"inline"`
	Class  int    `json:"class" bson:"class" validate:"required"`
	Stream string `json:"stream" bson:"stream" validate:"required"`
}
type Teachers struct {
	Users      `bson:"inline"`
	Department string  `json:"dept" bson:"dept" validate:"required,min=2,max=30"`
	Subject    string  `json:"subject" bson:"subject" validate:"required,min=2,max=30"`
	Experience float64 `json:"experience" bson:"experience" validate:"required"`
}

type User interface {
	Validate() error
	SetRole()
	EncryptPassword() error
	GetRolePssword() (string, string)
}

type Authenticate struct {
	Username string `json:"user" bson:"username,omitempty" validate:"required,min=4,max=50"`
	Password string `json:"password" bson:"password,omitempty" validate:"required,min=6,max=50"`
}

type ResponseStatus struct {
	Err     string `json:"error,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (u *Users) EncryptPassword() error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(password)
	if err != nil {
		return errors.New("unable to encrypt password")
	}
	return nil
}
func (s *Students) SetRole() {
	s.Uid = primitive.NewObjectID()
	s.Role = "student"
	s.IsValidated = false
	s.CreatedAt = time.Now()
}

func (t *Teachers) SetRole() {
	t.Uid = primitive.NewObjectID()
	t.Role = "teacher"
	t.IsValidated = true
	t.CreatedAt = time.Now()
}

var validate = validator.New()

func (s *Students) Validate() error {
	err := validate.Struct(s)
	if err != nil {
		fmt.Println(err)
		return errors.New("invalid user data")
	}
	return nil
}

func (t *Teachers) Validate() error {
	err := validate.Struct(t)
	if err != nil {
		return errors.New("invalid user data for teacher")
	}
	return nil
}

func (t *Teachers) GetRolePssword() (string, string) {
	return t.Password, t.Role
}

func (s *Students) GetRolePssword() (string, string) {
	return s.Password, s.Role
}
