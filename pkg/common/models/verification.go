package models

import "time"

type VerificationData struct {
	Email     *string    `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
	Code      *string    `json:"code,omitempty" bson:"code,omitempty" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}
