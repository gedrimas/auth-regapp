package model

import (

	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string            `json:"username"`
	Password      string            `json:"password"` 
	Email         string            `json:"email"` 
	Contacts      string			`json:"contacts"`
	User_type     string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Token         string            `json:"token"`
	User_id       string             `json:"user_id"`
	Refresh_token string            `json:"refresh_token"`
	Company_id 	  string			`json:"company_id"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}
