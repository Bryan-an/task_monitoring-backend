package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	Email  *bool `json:"email,omitempty" bson:"email,omitempty"`
	Mobile *bool `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

type Settings struct {
	Id            *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId        *string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Notifications *Notification       `json:"notifications,omitempty" bson:"notifications,omitempty"`
	Theme         *string             `json:"theme,omitempty" bson:"theme,omitempty"`
	CreatedAt     *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
