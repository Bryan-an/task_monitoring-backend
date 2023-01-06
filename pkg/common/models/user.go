package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      *string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     *string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  *string             `json:"password,omitempty" bson:"password,omitempty"`
	Role      *string             `json:"role,omitempty" bson:"role,omitempty"`
	Status    *string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UserDetails struct {
	ID    string
	Name  string
	Email string
}
