package models

import "time"

type VerificationData struct {
	Email     string    `json:"email" bson:"email,omitempty" binding:"required,email"`
	Code      string    `json:"code" bson:"code,omitempty" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at,omitempty"`
}
