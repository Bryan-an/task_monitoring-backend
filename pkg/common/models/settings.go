package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	Email  bool `json:"email" bson:"email"`
	Mobile bool `json:"mobile" bson:"mobile"`
}

type Settings struct {
	Id            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId        string             `json:"user_id" bson:"user_id,omitempty"`
	Notifications Notification       `json:"notifications" bson:"notifications,omitempty"`
	Security      string             `json:"security" bson:"security,omitempty"`
	Theme         string             `json:"theme" bson:"theme,omitempty"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}
