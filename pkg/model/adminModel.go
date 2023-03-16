package model

type Approve struct {
	Username string `json:"username" bson:"username" validate:"required"`
}
