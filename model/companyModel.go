package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID            primitive.ObjectID `bson:"_id"`
	Company       string            `json:"company"`
	Conpmay_id    string            `json:"company_id"`
	Admin_id      string            `json:"admin_id"`
}
