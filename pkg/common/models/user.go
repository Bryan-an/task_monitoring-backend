package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	Role      string             `json:"role" bson:"role,omitempty"`
	Status    string             `json:"status" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

type UserDetails struct {
	ID    string
	Name  string
	Email string
}
