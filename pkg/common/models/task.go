package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId      *string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Title       *string             `json:"title,omitempty" bson:"title,omitempty"`
	Description *string             `json:"description,omitempty" bson:"description,omitempty"`
	Labels      *[]string           `json:"labels,omitempty" bson:"labels,omitempty"`
	Priority    *string             `json:"priority,omitempty" bson:"priority,omitempty"`
	Complexity  *string             `json:"complexity,omitempty" bson:"complexity,omitempty"`
	Date        *time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	From        *time.Time          `json:"from,omitempty" bson:"from,omitempty"`
	To          *time.Time          `json:"to,omitempty" bson:"to,omitempty"`
	Done        *bool               `json:"done,omitempty" bson:"done,omitempty"`
	Remind      *bool               `json:"remind,omitempty" bson:"remind,omitempty"`
	Status      *string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
